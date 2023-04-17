package microtoken

import (
	"encoding/base64"
	"log"
	"testing"
)

func TestToken(t *testing.T) {
	var token = `REDACTED (todo: make this test sustainable)`

	codec := Codec{
		ClientID:   "c397d59039bbf1a2930ab1adff6c9e17",
		AuthDomain: "auth.tmpim.pw",
	}

	result, err := codec.EncodeToken([]byte(token))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	log.Println(base64.RawStdEncoding.EncodeToString(result))

	jwt, err := codec.DecodeToken(HS256Header, result)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	log.Println(string(jwt))
}
