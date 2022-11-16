package configstore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/scrypt"
)

var dataDir = "/data"

func writeConfig(kind string, config map[string]string, path string, audit string) error {
	db, err := bolt.Open(filepath.Join(dataDir, path), 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(kind))
		if err != nil {
			log.Fatal(err)
		}
		if err := b.Put([]byte("_aud"), []byte(audit)); err != nil {
			return err
		}
		if err := b.Put([]byte("_t    s"), []byte(time.Now().Format(time.RFC3339))); err != nil {
			return err
		}
		for key, value := range config {
			if strings.HasPrefix(key, "_") {
				continue
			}
			if err := b.Put([]byte(key), []byte(value)); err != nil {
				return err
			}
		}
		return nil
	})
}

func readConfig(kind, path string) (map[string]string, error) {
	config := make(map[string]string)
	db, err := bolt.Open(filepath.Join(dataDir, path), 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	return config, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			config[string(k)] = string(v)
			return nil
		})
	})
}

func encryptValues(passphrase string, config map[string]string) (map[string]string, error) {
	encryptedSecrets := map[string]string{}
	key, salt, err := deriveKey(passphrase, nil)
	if err != nil {
		return nil, err
	}
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	for name, value := range config {
		ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)
		ciphertext = append(ciphertext, salt...)
		encryptedSecrets[name] = base64.StdEncoding.EncodeToString(ciphertext)
	}
	return encryptedSecrets, nil
}

func decryptValues(passphrase string, config map[string]string) (map[string]string, error) {
	decryptedSecrets := map[string]string{}
	for name, value := range config {
		if strings.HasPrefix(name, "_") {
			continue
		}
		value, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			fmt.Println(value)
			return nil, err
		}
		salt, data := value[len(value)-32:], value[:len(value)-32]
		key, _, err := deriveKey(passphrase, salt)
		if err != nil {
			return nil, err
		}
		blockCipher, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		gcm, err := cipher.NewGCM(blockCipher)
		if err != nil {
			return nil, err
		}
		nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return nil, err
		}
		decryptedSecrets[name] = string(plaintext)
	}
	return decryptedSecrets, nil
}

func deriveKey(passphrase string, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key([]byte(passphrase), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
