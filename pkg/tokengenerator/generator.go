package tokengenerator

import (
	"crypto/rand"
	"encoding/base64"
)

type TokenGenerator struct{}

func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

func (g *TokenGenerator) GenerateRandomToken(len int) (string, error) {
	bytes := make([]byte, len)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(bytes)

	return token, nil
}
