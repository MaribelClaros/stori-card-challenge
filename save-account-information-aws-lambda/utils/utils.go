package utils

import (
	"math/rand"
	"time"
)

type UserIDGenerator struct {
}

func NewUserIDGenerator() *UserIDGenerator {
	return &UserIDGenerator{}
}

func (g *UserIDGenerator) GenerateID() int64 {
	rand.Seed(time.Now().UnixNano())

	// Generate a random 6-digit integer
	randomInt := rand.Intn(900000) + 100000

	// Convert the random integer to int64
	id := int64(randomInt)
	return id
}
