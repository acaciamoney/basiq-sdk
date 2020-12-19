package utilities

import (
	"log"
	"testing"
)

func TestGetToken(t *testing.T) {
	o, e := GetToken("123", "1.0")
	if e != nil {
		log.Print(e)
	}
	log.Print(o)
}
