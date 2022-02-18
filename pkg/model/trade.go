package model

import "time"

type Trade struct {
	//Price of the settled trade
	Price float64

	// Side of trade, may be Buy or Sell
	Side string

	// Size Amount of currency traded
	Size float64

	Timestamp time.Time

	// Market base-quote pair
	Market string
}
