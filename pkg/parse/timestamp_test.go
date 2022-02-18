package parse

import (
	"log"
	"testing"
	"time"
)

var (
	timezone, _ = time.LoadLocation("UTC")
)

func TestParseTimestamp(t *testing.T) {
	date, err := ParseTimestamp("2020-03-14T12:02:05.000234", timezone)
	if err != nil {
		t.Log(err)
	}
	log.Print(date)
	dateTwo, errTwo := ParseTimestamp("2020-03-14 12:02:05.000234", timezone)
	if errTwo != nil {
		t.Log(errTwo)
	}
	log.Print(dateTwo)
	dateThree, errThree := ParseTimestamp("2020-03-14T12:02:05.000234Z", timezone)
	if errThree != nil {
		t.Log(errThree)
	}
	log.Print(dateThree)
	_, errFour := ParseTimestamp("2020-0314 12:005.000234Z", timezone)
	if errFour == nil {
		t.Log("timestamp should not be parsed")
	}
}
