package wsgateway

import (
	"wsgateaway/exchanges"
)

type ExchangeWSClient interface {
	Connect() error
	Listen([]string, func(map[string]interface{})) error
	Close()
}

func BuildWSClient(exchange, APIKey, APISecret string) (ExchangeWSClient, error) {
	if exchange == "FTX" {
		return exchanges.NewFTXWSClient(APIKey, APISecret), nil
	}
	return nil, nil
}