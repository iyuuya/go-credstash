package cipher

import (
	"crypto/aes"
	cip "crypto/cipher"
)

type Cipher struct {
	block cip.Block
}

func NewCipher(key []byte) (*Cipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &Cipher{block}, nil
}

func (c *Cipher) Encrypt(value []byte) []byte {
	iv := make([]byte, aes.BlockSize)
	iv[len(iv)-1] = 1

	ctr := cip.NewCTR(c.block, iv)
	encrypted := make([]byte, len(value))
	ctr.XORKeyStream(encrypted, value)
	return encrypted
}

func (c *Cipher) Decrypt(value []byte) []byte {
	// IV (Initialization Vector) should be the same as used during encryption
	iv := make([]byte, aes.BlockSize)
	iv[len(iv)-1] = 1

	ctr := cip.NewCTR(c.block, iv)
	decrypted := make([]byte, len(value))
	ctr.XORKeyStream(decrypted, value)

	return decrypted
}
