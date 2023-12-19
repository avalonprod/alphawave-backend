package codegenerator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/xlzd/gotp"
)

type CodeGenerator struct {
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (g *CodeGenerator) GenerateUniqueCode() string {
	rand.Seed(time.Now().UnixNano())

	code := rand.Intn(90000) + 10000

	return fmt.Sprintf("%05d", code)
}

func (g *CodeGenerator) RandomSecret(len int) string {
	return gotp.RandomSecret(len)
}

func (g *CodeGenerator) GenerateUUID() string {
	uuid := uuid.New().String()
	return uuid
}
