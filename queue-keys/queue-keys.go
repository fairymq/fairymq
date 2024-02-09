package queue_keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
)

type PrivateKeyConfig struct {
	KeyDirectory string
}

// PrivateKeyContainer holds the keys used to decrypt queue messages
type PrivateKeyContainer struct {
	config PrivateKeyConfig
	keys   map[string]*rsa.PrivateKey
	mut    *sync.RWMutex
}

func NewPrivateKeyContainer(config PrivateKeyConfig) PrivateKeyContainer {
	return PrivateKeyContainer{
		config: config,
		keys:   make(map[string]*rsa.PrivateKey),
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

	pk.Add(queue, privateKey)

	return nil
}

func (pk *PrivateKeyContainer) Add(queue string, key *rsa.PrivateKey) {
	pk.mut.Lock()
	defer pk.mut.Unlock()

	for _, k := range pk.keys {
		if k.Equal(key) {
			return
		}
	}

	pk.keys[queue] = key
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

			queue := file.Name()[:strings.Index(file.Name(), ".private.pem")]
			pk.Add(queue, privateKey)
		}
	}

	pk.mut.RLock()
	defer pk.mut.RUnlock()

	return nil
}

func (pk *PrivateKeyContainer) DecryptMessage(buf []byte) (string, []byte, error) {
	var queue string
	var decrypted []byte
	var err error

	pk.mut.RLock()
	defer pk.mut.RUnlock()

	for q, key := range pk.keys {
		decrypted, err = rsa.DecryptPKCS1v15(rand.Reader, key, buf)
		if err == nil {
			queue = q
			break
		}
	}

	return queue, decrypted, err
}
