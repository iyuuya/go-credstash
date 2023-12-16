package cipher

import (
	"encoding/base64"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

var DefaultKMSKeyID = "alias/credstash"

type CipherKey struct {
	dataKey    []byte
	hmacKey    []byte
	wrappedKey []byte
}

func Generate(kmsClient *kms.Client, kmsKeyId *string, ctx map[string]string) (*CipherKey, error) {
	var keyId *string
	if kmsKeyId == nil {
		keyId = &DefaultKMSKeyID
	} else {
		keyId = kmsKeyId
	}
	out, err := kmsClient.GenerateDataKey(context.TODO(), &kms.GenerateDataKeyInput{
		KeyId:             keyId,
		NumberOfBytes:     aws.Int32(64),
		EncryptionContext: ctx,
	})
	if err != nil {
		return nil, err
	}

	return &CipherKey{
		dataKey:    out.Plaintext[0:32],
		hmacKey:    out.Plaintext[32:],
		wrappedKey: out.CiphertextBlob,
	}, nil
}

func Decrypt(wrappedKey []byte, kmsClient *kms.Client, ctx map[string]string) (*CipherKey, error) {
	out, err := kmsClient.Decrypt(context.TODO(), &kms.DecryptInput{
		CiphertextBlob:    wrappedKey,
		EncryptionContext: ctx,
	})
	if err != nil {
		return nil, err
	}

	return &CipherKey{
		dataKey:    out.Plaintext[0:32],
		hmacKey:    out.Plaintext[32:],
		wrappedKey: wrappedKey,
	}, nil
}

func (ck *CipherKey) HMAC(message []byte) string {
	h := hmac.New(sha256.New, ck.hmacKey)
	h.Write(message)
	s := h.Sum(nil)
	return hex.EncodeToString(s)
}

func (ck *CipherKey) Encrypt(message []byte) ([]byte, error) {
	c, err := NewCipher(ck.dataKey)
	if err != nil {
		return nil, err
	}
	return c.Encrypt(message), nil
}

func (ck *CipherKey) Decrypt(message []byte) ([]byte, error) {
	c, err := NewCipher(ck.dataKey)
	if err != nil {
		return nil, err
	}
	return c.Decrypt(message), nil
}

func (ck *CipherKey) Base64WrappedKey() string {
	return base64.StdEncoding.EncodeToString(ck.wrappedKey)
}
