package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/AlnsV/go-crypto-ws-gateway/wsgateway/internal"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const address = "wss://ftx.com/ws/"

type FTXWSClient struct {
	APIKey    string
	APISecret string

	ws *internal.WebsocketClient

	logger *logrus.Logger
}

func NewFTXWSClient(APIKey string, APISecret string) *FTXWSClient {
	return &FTXWSClient{
		APIKey:    APIKey,
		APISecret: APISecret,
		ws:        internal.NewWebsocketClient(address, true),
		logger:    logrus.New(),
	}
}

func signedSecret(msg, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)

	return hex.EncodeToString(mac.Sum(nil))
}

func (f FTXWSClient) Connect() error {
	response, err := f.ws.Connect(http.Header{})
	if err != nil {
		return err
	}
	f.logger.Infoln(response)

	ts := time.Now().Second()
	loginMsg := map[string]interface{}{
		"op": "login",
		"args": map[string]interface{}{
			"key": f.APIKey,
			"sign": signedSecret(
				[]byte(fmt.Sprintf("%dwebsocket_login", ts)),
				[]byte(f.APISecret),
			),
		},
	}

	f.logger.Infoln(loginMsg)

	err = f.ws.SendMessageJSON(loginMsg)
	if err != nil {
		return err
	}
	return nil
}

func (f FTXWSClient) subscribe(pairs []string) error {
	var err error
	for _, pair := range pairs {
		loginMsg := map[string]interface{}{
			"op":      "subscribe",
			"channel": "trades",
			"market":  pair,
		}

		f.logger.Infoln(loginMsg)

		err = f.ws.SendMessageJSON(loginMsg)
		if err != nil {
			break
		}
	}
	return err
}

func (f FTXWSClient) Listen(instruments []string, receiver func(map[string]interface{})) error {
	err := f.subscribe(instruments)
	if err != nil {
		return err
	}
	messageContainer := make(chan map[string]interface{})
	go f.ws.Listen(messageContainer)

	go func() {
		for {
			select {
			case msg := <-messageContainer:
				f.logger.Infoln(msg)
				// TODO(JV): Unify msg format
				receiver(msg)
			}
		}
	}()

	return nil
}

func (f FTXWSClient) Close() {
	//TODO implement me
	panic("implement me")
}
