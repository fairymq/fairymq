package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"sync"
	"time"
)

type SyncQueue struct {
	Name           string        `json:"Name"`
	ExpireMessages bool          `json:"ExpireMessages"`
	ExpiryTime     uint          `json:"ExpiryTime"`
	Messages       []SyncMessage `json:"Messages"`
}

type SyncMessage struct {
	Key       string    `json:"Key"`
	Data      []byte    `json:"Data"`
	Timestamp time.Time `json:"Timestamp"`
}

type NodeMeta struct {
	Name           string
	BindAddress    string
	MemberlistPort uint
}

type Delegate struct {
	fairyMQ *FairyMQ
}

func (delegate *Delegate) NodeMeta(limit int) []byte {
	meta := NodeMeta{
		Name:           fmt.Sprintf("%s:%d", delegate.fairyMQ.Config.BindAddress, delegate.fairyMQ.Config.MemberlistPort),
		BindAddress:    delegate.fairyMQ.Config.BindAddress,
		MemberlistPort: delegate.fairyMQ.Config.MemberlistPort,
	}
	mb := make([]byte, limit)
	mb, err := json.Marshal(meta)
	if err != nil {
		log.Println("Error: ", err)
		return []byte{}
	}
	return mb
}

func (delegate *Delegate) NotifyMsg([]byte) {
	// No-Op
}

func (delegate *Delegate) GetBroadcasts(overhead, limit int) [][]byte {
	// No-Op
	return [][]byte{}
}

func (delegate *Delegate) LocalState(join bool) []byte {
	var queues []SyncQueue

	for queueName, mut := range fairyMQ.QueueMutexes {
		mut.Lock()

		// Extract messages from queue
		var messages []SyncMessage
		for _, m := range fairyMQ.Queues[queueName].Messages {
			messages = append(messages, SyncMessage{
				Key:       m.Key,
				Data:      m.Data,
				Timestamp: m.Timestamp,
			})
		}

		queues = append(queues, SyncQueue{
			Name:           queueName,
			ExpireMessages: fairyMQ.Queues[queueName].ExpireMessages,
			ExpiryTime:     fairyMQ.Queues[queueName].ExpiryTime,
			Messages:       messages,
		})
		mut.Unlock()
	}

	b, err := json.Marshal(queues)
	if err != nil {
		log.Println("Could not encode state for sync: ", err.Error())
		return []byte{}
	}

	return b
}

func (delegate *Delegate) MergeRemoteState(buf []byte, join bool) {
	var queues []SyncQueue

	err := json.Unmarshal(buf, &queues)
	if err != nil {
		log.Println("Could not decode state for merge: ", err.Error())
		return
	}

	for _, q := range queues {
		var messages []Message
		mut, ok := fairyMQ.QueueMutexes[q.Name]

		if !ok { // If queue does not exist, add it.
			for _, m := range q.Messages {
				messages = append(messages, Message{
					Key:                   m.Key,
					Data:                  m.Data,
					Timestamp:             m.Timestamp,
					AcknowledgedConsumers: []Consumer{},
				})
			}

			fairyMQ.QueueMutexes[q.Name] = &sync.Mutex{}
			fairyMQ.QueueMutexes[q.Name].Lock()

			fairyMQ.Queues[q.Name] = &Queue{
				ExpireMessages: q.ExpireMessages,
				ExpiryTime:     q.ExpiryTime,
				Messages:       messages,
				Consumers:      []string{},
			}

			fairyMQ.QueueMutexes[q.Name].Unlock()
			continue
		}

		// If queue exists, merge the messages.
		mut.Lock()
		for _, m := range q.Messages {
			msgIdx := slices.IndexFunc(fairyMQ.Queues[q.Name].Messages, func(message Message) bool {
				return (m.Key == message.Key) && (m.Timestamp == message.Timestamp) && bytes.Equal(m.Data, message.Data)
			})
			if msgIdx == -1 {
				// Current message is not contained in the messages, add the message
				fairyMQ.Queues[q.Name].Messages = append(fairyMQ.Queues[q.Name].Messages, Message{
					Key:                   m.Key,
					Data:                  m.Data,
					Timestamp:             m.Timestamp,
					AcknowledgedConsumers: []Consumer{},
				})
			}
		}
		// Sort the messages by timestamp
		slices.SortFunc(fairyMQ.Queues[q.Name].Messages, func(a, b Message) int {
			switch {
			case a.Timestamp.Before(b.Timestamp):
				return 1
			case a.Timestamp.After(b.Timestamp):
				return -1
			default:
				return 0
			}
		})
		mut.Unlock()
	}
}
