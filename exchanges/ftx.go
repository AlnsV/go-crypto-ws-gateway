package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/AlnsV/go-crypto-ws-gateway/internal"
	"github.com/AlnsV/go-crypto-ws-gateway/pkg/model"
	"github.com/AlnsV/go-crypto-ws-gateway/pkg/parse"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const address = "wss://ftx.com/ws/"

var (
	timezone, _ = time.LoadLocation("UTC")
)

type FTXWSClient struct {
	APIKey string

	APISecret string

	ws *internal.WebsocketClient
}

func NewFTXWSClient(APIKey string, APISecret string) *FTXWSClient {
	return &FTXWSClient{
		APIKey:    APIKey,
		APISecret: APISecret,
		ws:        internal.NewWebsocketClient(address, true),
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
	logger.Infoln(response)

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

	logger.Infoln(loginMsg)

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

		logger.Infoln(loginMsg)

		err = f.ws.SendMessageJSON(loginMsg)
		if err != nil {
			break
		}
	}
	return err
}

func (f FTXWSClient) Listen(instruments []string, receiver func(trade *model.Trade)) error {
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
				logger.Infoln(msg)

				if trades, ok := msg["data"]; ok {
					for _, rawTrade := range trades.([]interface{}) {
						trade := rawTrade.(map[string]interface{})
						timestamp, _ := parse.ParseTimestamp(trade["time"].(string), timezone)

						newTrade := &model.Trade{
							Price:     trade["price"].(float64),
							Side:      trade["side"].(string),
							Size:      trade["size"].(float64),
							Timestamp: timestamp,
							Market:    msg["market"].(string),
						}
						receiver(newTrade)
					}
				}
			}
		}
	}()

	return nil
}

func (f FTXWSClient) Close() {
	f.ws.Close()
}
