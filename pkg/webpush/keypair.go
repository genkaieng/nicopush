package webpush

import (
	"crypto/ecdh"
	"crypto/rand"
	"log"
)

type KeyPair struct {
	Priv []byte
	Pub  []byte
}

func NewKeyPair() KeyPair {
	curve := ecdh.P256()

	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	return KeyPair{
		Priv: privateKey.Bytes(),
		Pub:  privateKey.PublicKey().Bytes(),
	}
}
