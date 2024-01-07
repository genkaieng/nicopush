package webpush

import (
	"crypto/rand"
	"log"
)

type Auth []byte

func NewAuth() Auth {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
