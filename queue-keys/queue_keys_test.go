package queue_keys

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"
	"testing"
)

const (
	KeyDirectory = "test/keys"
)

type Message struct {
	Raw       []byte
	Encrypted []byte
}

func TestDefaultPrivateKeyContainer(t *testing.T) {
	pkc := NewDefaultPrivateKeyContainer(PrivateKeyConfig{
		KeyDirectory: KeyDirectory,
	})

	// Generate the queue-key pairs and create the files
	queues := []string{"test-queue-one", "test-queue-two", "test.queue.three"}
	for _, queue := range queues {
		if err := pkc.GenerateQueueKeypair(queue); err != nil {
			t.Error(err)
		}
	}

	// Load the keys
	if err := pkc.LoadKeys(); err != nil {
		t.Error(err)
	}

	// Check all the expected keys for all the queues have been loaded into memory
	if len(pkc.keys) != len(queues) {
		t.Errorf("number of keys loaded %d is not equal to length of key map %d", len(pkc.keys), len(queues))
	}
	for _, queue := range queues {
		if _, ok := pkc.keys[queue]; !ok {
			t.Errorf("key for queue %s not found in loaded key map", queue)
		}
	}

	// Retrieve the public keys
	publicKeys := make(map[string]*rsa.PublicKey)
	for _, queue := range queues {
		publicKeyPEM, err := os.ReadFile(path.Join(KeyDirectory, fmt.Sprintf("%s.public.pem", queue)))
		if err != nil {
			t.Error(err)
		}
		publicKeyBlock, _ := pem.Decode(publicKeyPEM)
		publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
		if err != nil {
			t.Error(err)
		}
		publicKeys[queue] = publicKey.(*rsa.PublicKey)
	}

	// Encrypt some messages with each public key
	messages := make(map[string]Message)
	for _, queue := range queues {
		raw := []byte(fmt.Sprintf("Test message for queue %s", queue))
		encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKeys[queue], raw)
		if err != nil {
			t.Error(err)
		}
		messages[queue] = Message{
			Raw:       raw,
			Encrypted: encrypted,
		}
	}

	// Decrypt each of the messages encrypted in the previous step
	for queue, message := range messages {
		q, msg, err := pkc.DecryptMessage(message.Encrypted)
		if err != nil {
			t.Error(err)
		}
		if q != queue {
			t.Errorf("wrong queue name %s, expected %s", q, queue)
		}
		if !bytes.Equal(msg, message.Raw) {
			t.Error("decrypted message not equal to raw message")
		}
	}

	// Clean up the key directory
	if err := os.RemoveAll(KeyDirectory); err != nil {
		t.Error(err)
	}
}
