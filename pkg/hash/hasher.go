package hash

import (
	"crypto/sha1"
	"fmt"
)

type Hasher struct {
	passwordSolt string
}

func NewHasher(passwordSolt string) *Hasher {
	return &Hasher{
		passwordSolt: passwordSolt,
	}
}

func (h *Hasher) Hash(input string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(input)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.passwordSolt))), nil
}
