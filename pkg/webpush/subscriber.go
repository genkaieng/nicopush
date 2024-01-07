package webpush

import (
	"encoding/base64"
	"errors"
	"log"
	"regexp"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Subscriber struct {
	conn    *websocket.Conn
	KeyPair KeyPair
	Auth    Auth
}

func Connect(keyPair KeyPair, auth []byte) *Subscriber {
	conn, _, err := websocket.DefaultDialer.Dial("wss://push.services.mozilla.com/", nil)
	if err != nil {
		panic(err)
	}

	return &Subscriber{
		conn:    conn,
		KeyPair: keyPair,
		Auth:    auth,
	}
}

func (sub *Subscriber) SendHello(uaid string) (string, error) {
	var b []byte
	if uaid != "" {
		b = []byte(`{"messageType":"hello","use_webpush":true,"uaid":"` + uaid + `"}`)
	} else {
		b = []byte(`{"messageType":"hello","use_webpush":true}`)
	}
	err := sub.conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return "", err
	}
	_, b, err = sub.conn.ReadMessage()
	if err != nil {
		return "", err
	}
	msg := string(b)
	if checkMessageType(msg) != Hello {
		return "", errors.New("Unexpected messageType")
	}
	log.Println("INFO", msg)
	re, err := regexp.Compile(`"uaid":"(\w+)"`)
	if err != nil {
		return "", err
	}
	result := re.FindStringSubmatch(msg)
	return result[1], nil
}

func (sub *Subscriber) Subscribe(channelID string, notification chan string, pushEndpoint chan string) error {
	// FireFoxなどのブラウザでニコニコの設置画面開きブラウザの開発ツールを開く。
	// アプリケーションタブからsw.jsの内容を表示すると以下のキーが埋め込まれているのでこれを拝借。
	// https://account.nicovideo.jp/my/account
	key := base64.URLEncoding.EncodeToString([]byte{4, 45, 60, 21, 218, 246, 36, 40, 82, 47, 73, 43, 230, 41, 142, 247, 210, 250, 205, 145, 186, 70, 125, 45, 4, 5, 141, 78, 90, 217, 124, 155, 108, 14, 135, 128, 190, 98, 82, 107, 176, 167, 80, 225, 233, 54, 23, 121, 204, 233, 52, 98, 116, 83, 160, 67, 147, 227, 182, 11, 122, 223, 3, 166, 40})

	err := sub.conn.WriteMessage(websocket.TextMessage, []byte(`{"messageType":"register","channelID":"`+channelID+`","key":"`+key+`"}`))
	if err != nil {
		panic(err)
	}
	for {
		_, b, err := sub.conn.ReadMessage()
		if err != nil {
			return err
		}
		msg := string(b)
		switch checkMessageType(msg) {
		case Register:
			log.Println("INFO", msg)
			re, err := regexp.Compile(`"pushEndpoint":"(https?://[\w-./]+)"`)
			if err != nil {
				return err
			}
			result := re.FindStringSubmatch(msg)
			pushEndpoint <- result[1]
		case Notification:
			re, err := regexp.Compile(`"data":"([\w-/\+]+)"`)
			if err != nil {
				return err
			}
			result := re.FindStringSubmatch(msg)
			b, err := Decrypt(sub.KeyPair.Pub, sub.KeyPair.Priv, sub.Auth, result[1])
			if err != nil {
				log.Println("ERROR", err)
				continue
			}
			notification <- string(b)
		case Unregister:
			log.Println("INFO", msg)
			return nil
		default:
			log.Println("WARN", "Received unexpected messageType", msg)
		}
	}
}

type MessageType int

const (
	Unknown MessageType = iota
	Hello
	Register
	Unregister
	Notification
)

func checkMessageType(m string) MessageType {
	re, err := regexp.Compile(`"messageType":"(\w+)"`)
	if err != nil {
		panic(err)
	}
	result := re.FindStringSubmatch(m)
	switch result[1] {
	case "hello":
		return Hello
	case "register":
		return Register
	case "unregister":
		return Unregister
	case "notification":
		return Notification
	default:
		return Unknown
	}
}

func Decrypt(pub, priv, auth []byte, data string) ([]byte, error) {
	salt, senderKey, ciphertext, rs, err := decodeMessage(data)
	if err != nil {
		return nil, err
	}
	ikm, err := computeSharedKey(priv, senderKey)
	if err != nil {
		return nil, err
	}
	cek, nonce, err := deriveKeyAndNonce(ikm, auth, pub, senderKey, salt)
	if err != nil {
		return nil, err
	}
	result := decryptCipherText(ciphertext, cek, nonce, rs)
	return result, nil
}

func (sub *Subscriber) SendUnsubscribe(channelID string) error {
	return sub.conn.WriteMessage(websocket.TextMessage, []byte(`{"messageType":"unregister","channelID":"`+channelID+`"}`))
}

func (sub *Subscriber) Close() {
	sub.conn.Close()
}

func GenChannelID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}
