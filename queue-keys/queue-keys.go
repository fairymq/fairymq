package queue_keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"sync"
)

type PrivateKeyConfig struct {
	KeyDirectory string
}

// PrivateKeyContainer holds the keys used to decrypt queue messages
type PrivateKeyContainer struct {
	config PrivateKeyConfig
	keys   []*rsa.PrivateKey
	mut    *sync.RWMutex
}

func NewPrivateKeyContainer(config PrivateKeyConfig) PrivateKeyContainer {
	return PrivateKeyContainer{
		config: config,
		keys:   make([]*rsa.PrivateKey, 0),
		mut:    &sync.RWMutex{},
	}
}

// GenerateQueueKeypair creates a queue keypair
func (pk *PrivateKeyContainer) GenerateQueueKeypair(queue string) error {
	if _, err := os.Stat(pk.config.KeyDirectory); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(pk.config.KeyDirectory, 0777)
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

	err = os.WriteFile(fmt.Sprintf("%s/%s.private.pem", pk.config.KeyDirectory, queue), privateKeyPEM, 0644)
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
	err = os.WriteFile(fmt.Sprintf("%s/%s.public.pem", pk.config.KeyDirectory, queue), publicKeyPEM, 0644)
	if err != nil {
		return err
	}

	pk.Add(privateKey)

	return nil
}

func (pk *PrivateKeyContainer) Add(key *rsa.PrivateKey) {
	pk.mut.Lock()
	defer pk.mut.Unlock()

	// TODO: Check if the key already exists in the key slice

	pk.keys = append(pk.keys, key)
}

func (pk *PrivateKeyContainer) LoadKeys() {
	// TODO: Load keys from files into memory
}
