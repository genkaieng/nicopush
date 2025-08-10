package cmd

import (
	"encoding/base64"

	"github.com/genkaieng/nicopush/pkg/webpush"
)

func Genkeys(args []string) int {
	keyPair := webpush.NewKeyPair()
	println("NICOPUSH_PUBLIC_KEY=" + base64.StdEncoding.EncodeToString(keyPair.Pub))
	println("NICOPUSH_PRIVATE_KEY=" + base64.StdEncoding.EncodeToString(keyPair.Priv))

	auth := webpush.NewAuth()
	println("NICOPUSH_AUTH=" + base64.StdEncoding.EncodeToString(auth))

	channelID := webpush.GenChannelID()
	println("NICOPUSH_CHANNEL_ID=" + channelID)

	return 0
}
