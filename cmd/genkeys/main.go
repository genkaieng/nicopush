package main

import (
	"encoding/base64"

	"github.com/genkaieng/nicopush-subscriber/pkg/webpush"
)

func main() {
	keyPair := webpush.NewKeyPair()
	println("PUBLIC_KEY=" + base64.StdEncoding.EncodeToString(keyPair.Pub))
	println("PRIVATE_KEY=" + base64.StdEncoding.EncodeToString(keyPair.Priv))

	auth := webpush.NewAuth()
	println("AUTH=" + base64.StdEncoding.EncodeToString(auth))

	channelID := webpush.GenChannelID()
	println("CHANNEL_ID=" + channelID)
}
