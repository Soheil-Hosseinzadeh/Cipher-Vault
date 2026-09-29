package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

const (
	memory     uint32 = 64 * 1024
	iterations uint32 = 3
	threads    uint8  = 2
	saltLength uint32 = 16
	keyLength  uint32 = 32
)

type Encrypter struct {
	PassKey string
}

func (e *Encrypter) Encrypt(plainText []byte) ([]byte, error) {

	// creating a random Salt for Argon2
	salt, err := generateRandomBytes(saltLength)
	if err != nil {
		return nil, err
	}

	// creating a key from user passKey by argon2
	key := argon2.IDKey([]byte(e.PassKey), salt, iterations, memory, threads, keyLength)

	// using AES-256
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//using GCM encryption method
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//initializing nonce
	nonce, err := generateRandomBytes(uint32(gcm.NonceSize()))
	if err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	cipherText = append(salt, cipherText...)

	return cipherText, nil
}

func (e *Encrypter) Decrypt(encryptedText []byte) ([]byte, error) {
	salt := encryptedText[:16]
	nonce := encryptedText[16:28]
	encryptedText = encryptedText[28:]

	// creating the same key from user passKey and the salt by argon2
	key := argon2.IDKey([]byte(e.PassKey), salt, iterations, memory, threads, keyLength)

	// using AES-256
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//using GCM encryption method
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	result, err := gcm.Open(nil, nonce, encryptedText, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
