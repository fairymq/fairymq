package queue_keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
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

	// TODO: Remove log statements

	if !slices.ContainsFunc(pk.keys, func(k *rsa.PrivateKey) bool {
		return key.Equal(k)
	}) {
		fmt.Println("Adding key to key list")
		pk.keys = append(pk.keys, key)
	} else {
		fmt.Println("This key already exists, skipping")
	}
}

func (pk *PrivateKeyContainer) LoadKeys() error {
	files, err := os.ReadDir(pk.config.KeyDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "private.pem") {
			privateKeyPEM, err := os.ReadFile(path.Join(pk.config.KeyDirectory, file.Name()))
			if err != nil {
				return err
			}
			privateKeyBlock, _ := pem.Decode(privateKeyPEM)
			privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
			if err != nil {
				return err
			}
			pk.Add(privateKey)
		}
	}

	return nil
}
