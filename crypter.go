package linkpearl

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Crypter is special object that handles data encryption and decryption
// apart from Matrix' standard encryption.
// It can encrypt and decrypt arbitrary data using secret key (password)
type Crypter struct {
	cipher    cipher.AEAD
	nonceSize int
}

// ErrInvalidData returned in provided encrypted data (ciphertext) is invalid
var ErrInvalidData = errors.New("invalid data")

// NewCrypter creates new Crypter
func NewCrypter(secretkey string) (*Crypter, error) {
	secret := []byte(secretkey)
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Crypter{
		cipher:    aesGCM,
		nonceSize: aesGCM.NonceSize(),
	}, nil
}

// Decrypt data
func (c *Crypter) Decrypt(data string) (string, error) {
	if len(data) < c.nonceSize {
		return "", ErrInvalidData
	}

	nonce := data[:c.nonceSize]
	ciphertext := data[c.nonceSize:]

	plaintext, err := c.cipher.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Encrypt data
func (c *Crypter) Encrypt(data string) (string, error) {
	nonce := make([]byte, c.nonceSize)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	encrypted := c.cipher.Seal(nonce, nonce, []byte(data), nil)
	return string(encrypted), nil
}
