package secret

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/iyuuya/go-credstash/cipher"
	"github.com/iyuuya/go-credstash/item"
)

type Secret struct {
	name           string
	value          string
	key            *cipher.CipherKey
	encryptedValue []byte
	hmac           []byte
	context        map[string]string
}

func NewSecret(
	name string,
	value string,
	key *cipher.CipherKey,
	encryptedValue []byte,
	hmac []byte,
	ctx map[string]string,
) *Secret {
	return &Secret{name, value, key, encryptedValue, hmac, ctx}
}

func Find(db *item.DynamoDB, kms *kms.Client, name string, ctx map[string]string, version string) (*Secret, error) {
	item, err := db.Get(name, version)
	if err != nil {
		return nil, err
	}

	enc, err := base64.StdEncoding.DecodeString(*item.GetContents())
	if err != nil {
		return nil, err
	}

	dec, err := base64.StdEncoding.DecodeString(item.GetKey())
	if err != nil {
		return nil, err
	}

	key, err := cipher.Decrypt(dec, kms, ctx)
	if err != nil {
		return nil, err
	}

	var h []byte
	if item.GetHMAC() != nil {
		h = []byte(*item.GetHMAC())
	}

	return &Secret{
		name:           name,
		key:            key,
		encryptedValue: enc,
		hmac:           h,
	}, nil
}

func (s *Secret) Encrypt(c *kms.Client, kmsKeyId *string, ctx map[string]string) error {
	key, err := cipher.Generate(c, kmsKeyId, ctx)
	if err != nil {
		return err
	}
	s.key = key

	enc, err := key.Encrypt([]byte(s.value))
	if err != nil {
		return err
	}
	s.encryptedValue = enc

	s.hmac = []byte(key.HMAC(s.encryptedValue))
	return nil
}

func (s *Secret) Save(db *item.DynamoDB) error {
	return db.Put(s.toItem(db))
}

func (s *Secret) IsFalsified() bool {
	return s.key.HMAC(s.encryptedValue) == string(s.hmac)
}

func (s *Secret) DecryptedValue() (string, error) {
	b, err := s.key.Decrypt(s.encryptedValue)
	return string(b), err
}

func (s *Secret) toItem(db *item.DynamoDB) *item.Item {
	key := s.key.Base64WrappedKey()
	contents := base64.StdEncoding.EncodeToString(s.encryptedValue)
	version := fmt.Sprintf("%019d", s.currentVersion(db)+1)
	hmac := string(s.hmac)

	return item.NewItem(
		key,
		&contents,
		&s.name,
		&version,
		&hmac,
	)
}

func (s *Secret) currentVersion(db *item.DynamoDB) int {
	items, err := db.Select(s.name, 1, "")
	if err != nil {
		return 0
	}
	if len(items) == 0 {
		return 0
	}

	v := *items[0].GetVersion()
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}

	return n
}
