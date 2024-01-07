package niconico

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ApiClient struct {
	Endpoint string
	Auth     []byte
	P256dh   []byte
	Session  string
}

func (c *ApiClient) Register() (int, []byte, error) {
	url := "https://api.push.nicovideo.jp/v1/nicopush/webpush/endpoints.json"

	data, err := NewRegisterMessage(c.Endpoint, c.Auth, c.P256dh)
	if err != nil {
		return 0, nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return 0, nil, err
	}

	setRequestHeader(&req.Header, c.Session)

	statusCode, b, err := send(req)
	if err != nil {
		return 0, nil, err
	}
	log.Println("INFO", "register nicopush", statusCode, string(b))
	return statusCode, b, nil
}

func (c *ApiClient) Unregister() (int, []byte, error) {
	url := "https://api.push.nicovideo.jp/v1/nicopush/webpush/endpoints.json"

	data, err := NewUnregisterMessage(c.Endpoint)
	if err != nil {
		return 0, nil, err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(data))
	if err != nil {
		return 0, nil, err
	}

	setRequestHeader(&req.Header, c.Session)

	statusCode, b, err := send(req)
	if err != nil {
		return 0, nil, err
	}
	log.Println("INFO", "unregister nicopush", statusCode, string(b))
	return statusCode, b, nil
}

func setRequestHeader(header *http.Header, session string) {
	header.Set("Accept", "*/*")
	header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0")
	header.Set("Cookie", "user_session="+session)
	header.Set("Origin", "https://account.nicovideo.jp")
	header.Set("Referer", "https://account.nicovideo.jp/")
	header.Set("X-Frontend-Id", "8")
	header.Set("X-Request-With", "https://account.nicovideo.jp/my/account?cmnhd_ref=device%3Dpc%26site%3Dniconico%26pos%3Duserpanel%26page%3Dtop")
}

func send(req *http.Request) (int, []byte, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}
	return res.StatusCode, body, nil
}

func NewRegisterMessage(endpoint string, auth, p256dh []byte) ([]byte, error) {
	return json.Marshal(Register{
		DestApp: "nico_account_webpush",
		Endpoint: EndPoint{
			Auth:     base64.StdEncoding.EncodeToString(auth),
			Endpoint: endpoint,
			P256dh:   base64.StdEncoding.EncodeToString(p256dh),
		},
	})
}

func NewUnregisterMessage(endpoint string) ([]byte, error) {
	return json.Marshal(Unregister{
		DestApp: "nico_account_webpush",
		Endpoint: EndPoint{
			Endpoint: endpoint,
		},
	})
}

type Register struct {
	DestApp  string   `json:"destApp"`
	Endpoint EndPoint `json:"endpoint"`
}

type Unregister struct {
	DestApp  string   `json:"destApp"`
	Endpoint EndPoint `json:"endpoint"`
}

type EndPoint struct {
	Auth     string `json:"auth,omitempty"`
	Endpoint string `json:"endpoint"`
	P256dh   string `json:"p256dh,omitempty"`
}
