package cipher_test

import (
	"bytes"
	"testing"

	"github.com/iyuuya/go-credstash/cipher"
)

func Test_Cipher(t *testing.T) {
	key := "00000000000000000000000000000000"
	input := "value"

	c, err := cipher.NewCipher([]byte(key))
	if err != nil {
		t.Error(err)
	}

	want := []byte{103, 57, 28, 87, 216}
	enc := c.Encrypt([]byte(input))
	if !bytes.Equal(want, enc) {
		t.Errorf("Encrypt() = %v, want %v", enc, want)
	}

	dec := c.Decrypt(enc)
	if input != string(dec) {
		t.Errorf("Decrypt() = %v, want %v", dec, []byte(input))
	}
}
