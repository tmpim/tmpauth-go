package main

import (
	"encoding/base64"
	"fmt"
	"syscall/js"

	"github.com/tmpim/tmpauth-go/microtoken"
)

func encodeMicrotoken() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return js.ValueOf("Invalid number of arguments, expected client ID and JWT")
		}

		codec := microtoken.Codec{
			ClientID:   args[0].String(),
			AuthDomain: "auth.tmpim.pw",
		}
		result, err := codec.EncodeToken([]byte(args[1].String()))
		if err != nil {
			return js.ValueOf(fmt.Sprintf("There was an error encoding your microtoken:\n%+v", err))
		}

		return js.ValueOf(base64.RawURLEncoding.EncodeToString(result))
	})
}

func main() {
	c := make(chan struct{})
	js.Global().Set("encodeMicrotoken", encodeMicrotoken())
	<-c
}
