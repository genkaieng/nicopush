package cmd

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/genkaieng/nicopush/pkg/niconico"
	"github.com/genkaieng/nicopush/pkg/webpush"
	_ "github.com/joho/godotenv/autoload"
)

func Subscribe(args []string) int {
	pub := os.Getenv("NICOPUSH_PUBLIC_KEY")
	priv := os.Getenv("NICOPUSH_PRIVATE_KEY")
	auth := os.Getenv("NICOPUSH_AUTH")
	channelID := os.Getenv("NICOPUSH_CHANNEL_ID")
	uaid := os.Getenv("NICOPUSH_UAID")
	session := os.Getenv("NICOLIVE_SESSION")
	if pub == "" || priv == "" || auth == "" || channelID == "" || session == "" {
		panic("環境変数を設定してください。")
	}
	keyPair := webpush.KeyPair{
		Pub:  base64ToByte(pub),
		Priv: base64ToByte(priv),
	}

	var sub *webpush.Subscriber
	var nicoApi *niconico.ApiClient

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-interrupt
		// Push通知のsubscribe終了処理

		if nicoApi != nil {
			_, _, err := nicoApi.Unregister()
			if err != nil {
				log.Println("ERROR", err)
			}
			nicoApi = nil
		}
		if sub != nil {
			sub.SendUnsubscribe(channelID)
		}
	}()

	sub = webpush.Connect(keyPair, base64ToByte(auth))
	defer sub.Close()

	newUaid, err := sub.SendHello(uaid)
	if err != nil {
		panic(err)
	}
	if uaid != newUaid {
		log.Println("WARN", "UAID="+newUaid)
	}

	notification := make(chan string)
	pushEndpoint := make(chan string)

	go func() {
		for {
			endpoint := <-pushEndpoint
			if nicoApi != nil {
				nicoApi.Unregister()
			}
			nicoApi = &niconico.ApiClient{
				Endpoint: endpoint,
				Auth:     sub.Auth,
				P256dh:   sub.KeyPair.Pub,
				Session:  session,
			}
			_, _, err := nicoApi.Register()
			if err != nil {
				panic(err)
			}
		}
	}()
	defer func() {
		if nicoApi != nil {
			nicoApi.Unregister()
			nicoApi = nil
		}
	}()

	go func() {
		for {
			n := <-notification
			fmt.Println("Notification:", n)
		}
	}()

	if err = sub.Subscribe(channelID, notification, pushEndpoint); err != nil {
		log.Println("ERROR", err)
	}

	return 0
}

func base64ToByte(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
