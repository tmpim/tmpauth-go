package main

import (
	"encoding/base64"
	"fmt"
	"strings"
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

		inputJWT := strings.TrimSpace(args[1].String())

		result, err := codec.EncodeToken([]byte(inputJWT))
		if err != nil {
			return js.ValueOf(fmt.Sprintf("There was an error encoding your microtoken:\n%+v", err))
		}

		// verify the resulting token
		output, err := codec.DecodeToken(microtoken.HS256Header, result)
		if err != nil {
			return js.ValueOf(fmt.Sprintf("There was an error verifying your microtoken, "+
				"this is a bug and should be reported:\n%+v", err))
		}

		if string(output) != inputJWT {
			return js.ValueOf(fmt.Sprintf("Your microtoken failed verification, " +
				"this is a bug and should be reported"))
		}

		return js.ValueOf(base64.RawURLEncoding.EncodeToString(result))
	})
}

func main() {
	c := make(chan struct{})
	js.Global().Set("encodeMicrotoken", encodeMicrotoken())
	<-c
}
