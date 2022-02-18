package wsgateway_test

import (
	wsgateway "github.com/AlnsV/go-crypto-ws-gateway"
	"github.com/AlnsV/go-crypto-ws-gateway/pkg/model"
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"testing"
	"time"
)

type config struct {
	FTXAPIKey    string `env:"FTX_API_KEY"`
	FTXAPISecret string `env:"FTX_API_SECRET"`
}

func New() (*config, error) {
	var cfg config

	if err := env.Parse(&cfg); err != nil {
		return nil, errors.Wrap(err, "error with initializing config")
	}

	return &cfg, nil
}

func TestGateway(t *testing.T) {
	cfg, _ := New()
	client, err := wsgateway.BuildWSClient(
		"FTX",
		cfg.FTXAPIKey,
		cfg.FTXAPISecret,
	)
	if err != nil {
		logger.Error(err)
	}

	err = client.Connect()
	if err != nil {
		logger.Error(err)
	}

	err = client.Listen(
		[]string{"BTC-PERP", "SOL-PERP"},
		func(trade *model.Trade) {
			logger.Infoln(trade)
		},
	)
	if err != nil {
		logger.Error(err)
	}

	time.Sleep(2 * time.Second)
	client.Close()
}
