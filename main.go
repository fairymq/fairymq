/*
* fairyMQ
* Core
* ******************************************************************
* Originally authored by Alex Gaetano Padula
* Copyright (C) fairyMQ
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// FairyMQ is the fairyMQ system structure
type FairyMQ struct {
	UDPAddr       *net.UDPAddr       // UDP address representation
	Conn          *net.UDPConn       // Conn is the implementation of the Conn and PacketConn interfaces for UDP network connections
	Wg            *sync.WaitGroup    // Waitgroup pointer
	SignalChannel chan os.Signal     // Signal channel
	Queues        map[string]*Queue  // In-memory queues
	Consumers     []Consumer         // Consumer
	ContextCancel context.CancelFunc // To cancel on signal
	Context       context.Context    // For signal cancellation
}

// Queue is the fairyMQ queue structure
type Queue struct {
	ExpireMessages bool      // Expire messages and delete from queue
	ExpiryTime     uint      // Expiry in seconds; Default is 7200 (2 hours)
	Messages       []Message // Queue messages
	Consumers      []string  // Consumer addresses
}

// Consumer is a queue consumer
type Consumer struct {
	Queue   string // Name of queue
	Address string // Consumer address i.e 0.0.0.0:5992
}

// Message is a queue message
type Message struct {
	Data                  []byte     // Message data
	Timestamp             time.Time  // Message timestamp
	AcknowledgedConsumers []Consumer // Which consumers acknowledged this message? if any
}

// Global variables
var (
	fairyMQ *FairyMQ // Main fairyMQ pointer
)

func main() {
	fairyMQ = &FairyMQ{
		Wg:            &sync.WaitGroup{},       // Setting waitgroup pointer to hold go routines
		SignalChannel: make(chan os.Signal, 1), // Make signal channel
		Queues:        make(map[string]*Queue), // Make queues in-memory hashmap
	} // Set fairyMQ global pointer

	generateQueueKeypair := "" // If provided a new keypair will be created.
	flag.StringVar(&generateQueueKeypair, "generate-queue-key-pair", generateQueueKeypair, "Generates a new queue keypair.   --generate-queue-key-pair=yourqueuename")
	flag.Parse() // Parse flags if any

	// If queue provided generate a new keypair
	if generateQueueKeypair != "" {
		err := fairyMQ.GenerateQueueKeypair(generateQueueKeypair)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	fairyMQ.Context, fairyMQ.ContextCancel = context.WithCancel(context.Background()) // Set core context to cancel on signal

	signal.Notify(fairyMQ.SignalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGABRT) // Populate signal channel on signal

	fairyMQ.Wg.Add(1)
	go fairyMQ.SignalListener() // Start signal listener

	fairyMQ.Wg.Add(1)
	go fairyMQ.StartUDPListener() // Start UDP listener on default port 5991

	fairyMQ.RecoverQueue()

	fairyMQ.Wg.Wait() // Wait for all go routines

}

// SignalListener listens for system signals and gracefully shutsdown
func (fairyMQ *FairyMQ) SignalListener() {
	defer fairyMQ.Wg.Done()
	for {
		select {
		case sig := <-fairyMQ.SignalChannel:
			log.Println("received", sig)
			fairyMQ.ContextCancel()
			fairyMQ.Conn.Close()
			fairyMQ.Snapshot()
			return
		default:
			time.Sleep(time.Nanosecond * 10000)
			continue
		}
	}
}

// GenerateQueueKeypair creates a queue keypair
func (fairyMQ *FairyMQ) GenerateQueueKeypair(queue string) error {
	if _, err := os.Stat("./keys"); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("keys", 0777)
			if err != nil {
				return err
			}
		}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	err = os.WriteFile(fmt.Sprintf("keys/%s.private.pem", queue), privateKeyPEM, 0644)
	if err != nil {
		return err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	err = os.WriteFile(fmt.Sprintf("keys/%s.public.pem", queue), publicKeyPEM, 0644)
	if err != nil {
		return err
	}

	return nil
}

// SendToConsumers sends message to consumers of a queue
func (fairyMQ *FairyMQ) SendToConsumers(queue string, data []byte, message *Message) {

	for _, c := range fairyMQ.Consumers {
		if c.Queue == queue {
			attempts := 0 // Max attempts to reach server is 10

			// Resolve UDP address
			udpAddr, err := net.ResolveUDPAddr("udp", c.Address)
			if err != nil {
				continue
			}

			// Dial address
			conn, err := net.DialUDP("udp", nil, udpAddr)
			if err != nil {
				continue
			}

			publicKeyPEM, err := os.ReadFile(fmt.Sprintf("keys/%s.public.pem", queue))
			if err != nil {
				continue
			}

			publicKeyBlock, _ := pem.Decode(publicKeyPEM)
			publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
			if err != nil {
				continue
			}

			ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), data)
			if err != nil {
				continue
			}

			// Attempt consumer
			goto try

		try:

			// Send to server
			_, err = conn.Write(ciphertext)
			if err != nil {
				continue
			}

			// If nothing received in 60 milliseconds.  Retry
			err = conn.SetReadDeadline(time.Now().Add(60 * time.Millisecond))
			if err != nil {
				continue
			}

			// Read from consumer
			res, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					attempts += 1

					if attempts < 10 {
						goto try
					} else {
						continue
					}
				} else {
					continue
				}
			}

			if strings.HasPrefix(res, "ACK") {
				message.AcknowledgedConsumers = append(message.AcknowledgedConsumers, c)
			}

		}
	}
}

// RecoverQueue recovers latest persisted queue
func (fairyMQ *FairyMQ) RecoverQueue() {
	if _, err := os.Stat("snapshots"); os.IsNotExist(err) {
		return
	}

	snapshots, err := os.ReadDir("snapshots")
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(snapshots, func(i, j int) bool {
		fileI, err := snapshots[i].Info()
		if err != nil {
			return false
		}
		fileJ, err := snapshots[j].Info()
		if err != nil {
			return false
		}
		return fileI.ModTime().After(fileJ.ModTime())
	})

	for _, snapshot := range snapshots {
		snapshotFile, err := os.Open(fmt.Sprintf("snapshots/%s", snapshot.Name()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dataDecoder := gob.NewDecoder(snapshotFile)
		err = dataDecoder.Decode(&fairyMQ.Queues)
		log.Println("Recovered from snapshot")
		break
	}

}

// Snapshot takes a snapshot of current queue
func (fairyMQ *FairyMQ) Snapshot() {
	if _, err := os.Stat("snapshots"); os.IsNotExist(err) {
		err := os.Mkdir("snapshots", 0777)
		if err != nil {
			log.Println("ERROR:", err.Error())
			return
		}
	}

	snapshot, err := os.Create(fmt.Sprintf("snapshots/queue.%d.snapshot", time.Now().Unix()))
	if err != nil {
		log.Println("ERROR:", err.Error())
		return
	}

	// serialize the data
	dataEncoder := gob.NewEncoder(snapshot)
	dataEncoder.Encode(fairyMQ.Queues)
}

// StartUDPListener starts listening and handling UDP connections
func (fairyMQ *FairyMQ) StartUDPListener() {
	defer fairyMQ.Wg.Done()
	var err error

	fairyMQ.UDPAddr, err = net.ResolveUDPAddr("udp", "0.0.0.0:5991")
	if err != nil {
		log.Println("ERROR: ", err.Error())
		fairyMQ.SignalChannel <- os.Interrupt
	}

	// Start listening for UDP packages on the given address
	fairyMQ.Conn, err = net.ListenUDP("udp", fairyMQ.UDPAddr)
	if err != nil {
		log.Println("ERROR: ", err.Error())
		fairyMQ.SignalChannel <- os.Interrupt
	}

	if _, err := os.Stat("./keys"); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("keys", 0777)
			if err != nil {
				log.Println("ERROR: ", err.Error())
				fairyMQ.SignalChannel <- os.Interrupt
			}
		}
	}

	for {
		fairyMQ.Conn.SetReadDeadline(time.Now().Add(time.Nanosecond * 10000)) // essentially keep listening until the client closes connection or cluster shuts down

		var buf [5120]byte

		n, addr, err := fairyMQ.Conn.ReadFromUDP(buf[0:])
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				if fairyMQ.Context.Err() != nil { // if signaled to shutdown
					break
				}
				continue
			} else {
				break
			}
		}

		keys, err := os.ReadDir("keys")
		if err != nil {
			continue
		}

		for _, key := range keys {
			if !key.IsDir() {

				if strings.HasSuffix(key.Name(), "private.pem") {
					privateKeyPEM, err := os.ReadFile("keys/" + key.Name())
					if err != nil {
						goto nack
					}
					privateKeyBlock, _ := pem.Decode(privateKeyPEM)
					privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
					if err != nil {
						goto nack
					}
					plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, buf[0:n])
					if err != nil {
						goto nack
					}

					queue := strings.Split(key.Name(), ".")[0]

					_, ok := fairyMQ.Queues[queue]
					if !ok {
						fairyMQ.Queues[queue] = &Queue{
							ExpireMessages: false,
							ExpiryTime:     7200,
							Messages:       []Message{},
							Consumers:      []string{},
						}
					}

					switch {
					case bytes.HasPrefix(plaintext, []byte("FIRST IN")):
						fairyMQ.Conn.WriteToUDP(append(fairyMQ.Queues[queue].Messages[0].Data, []byte("\r\n")...), addr)
						goto cont
					case bytes.HasPrefix(plaintext, []byte("LAST IN")):
						fairyMQ.Conn.WriteToUDP(append(fairyMQ.Queues[queue].Messages[len(fairyMQ.Queues[string(queue)].Messages)-1].Data, []byte("\r\n")...), addr)
						goto cont
					case bytes.HasPrefix(plaintext, []byte("LENGTH")):
						fairyMQ.Conn.WriteToUDP(append([]byte(fmt.Sprintf("%d messages", len(fairyMQ.Queues[string(queue)].Messages))), []byte("\r\n")...), addr)
						goto cont
					case bytes.HasPrefix(plaintext, []byte("POP")):
						fairyMQ.Queues[queue].Messages = fairyMQ.Queues[queue].Messages[:len(fairyMQ.Queues[string(queue)].Messages)-1]
						fairyMQ.Conn.WriteToUDP(append([]byte(fmt.Sprintf("ACK")), []byte("\r\n")...), addr)
						goto cont
					case bytes.HasPrefix(plaintext, []byte("SHIFT")):
						fairyMQ.Queues[queue].Messages = fairyMQ.Queues[queue].Messages[1:]
						fairyMQ.Conn.WriteToUDP(append([]byte(fmt.Sprintf("ACK")), []byte("\r\n")...), addr)
						goto cont
					case bytes.HasPrefix(plaintext, []byte("CLEAR")):
						delete(fairyMQ.Queues, queue)
						fairyMQ.Conn.WriteToUDP(append([]byte(fmt.Sprintf("ACK")), []byte("\r\n")...), addr)
						goto cont
					case bytes.HasPrefix(plaintext, []byte("NEW CONSUMER ")):
						spl := bytes.Split(plaintext, []byte("NEW CONSUMER "))

						for _, c := range fairyMQ.Consumers {
							if c.Queue == queue {
								if c.Address == strings.TrimSpace(string(spl[1])) {
									fairyMQ.Conn.WriteToUDP(append([]byte(fmt.Sprintf("NACK")), []byte("\r\n")...), addr)
									goto cont
								}
							}
						}

						fairyMQ.Consumers = append(fairyMQ.Consumers, Consumer{
							Queue:   queue,
							Address: strings.TrimSpace(string(spl[1])),
						})
						fairyMQ.Conn.WriteToUDP(append([]byte(fmt.Sprintf("ACK")), []byte("\r\n")...), addr)
						goto cont
					case bytes.HasPrefix(plaintext, []byte("REM CONSUMER ")):
						spl := bytes.Split(plaintext, []byte("REM CONSUMER "))
						fairyMQ.Consumers = append(fairyMQ.Consumers, Consumer{
							Queue:   queue,
							Address: strings.TrimSpace(string(spl[1])),
						})
						fairyMQ.Conn.WriteToUDP(append([]byte(fmt.Sprintf("ACK")), []byte("\r\n")...), addr)
						goto cont
					case bytes.HasPrefix(plaintext, []byte("LIST CONSUMERS")):
						var consumers []string

						for _, c := range fairyMQ.Consumers {
							if c.Queue == queue {
								consumers = append(consumers, c.Address)
							}
						}

						fairyMQ.Conn.WriteToUDP(append([]byte(fmt.Sprintf(strings.Join(consumers, ","))), []byte("\r\n")...), addr)
						goto cont
					case bytes.HasPrefix(plaintext, []byte("ENQUEUE")):
						spl := bytes.Split(plaintext, []byte("\r\n"))
						timestamp, err := strconv.ParseInt(string(spl[1]), 10, 64)
						if err != nil {
							goto cont
						}

						message := Message{
							Data:      spl[2],
							Timestamp: time.UnixMicro(timestamp),
						}

						fairyMQ.Queues[queue].Messages = append(fairyMQ.Queues[queue].Messages, message)

						go fairyMQ.SendToConsumers(queue, plaintext, &message)

						sort.Slice(fairyMQ.Queues[queue].Messages, func(i, j int) bool {
							return fairyMQ.Queues[queue].Messages[i].Timestamp.After(fairyMQ.Queues[queue].Messages[j].Timestamp)
						})

						fairyMQ.Conn.WriteToUDP([]byte("ACK\r\n"), addr)
						goto cont
					default:
						fairyMQ.Conn.WriteToUDP([]byte("NACK\r\n"), addr)
						goto cont
					}
				}
			}

		}

		fairyMQ.Conn.WriteToUDP([]byte("NACK\r\n"), addr)

	cont:
		continue

	nack:
		fairyMQ.Conn.WriteToUDP([]byte("NACK\r\n"), addr)

	}
}
