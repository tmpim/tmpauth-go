package main

import (
	"log"
	"net/http"

	"github.com/tmpim/tmpauth-go"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	ta, err := tmpauth.NewMini(tmpauth.MiniConfig{
		PublicKey:      "BN/PHEYgs0meH878gqpWl81WD3zEJ+ubih3RVYwFxaYXxHF+5tgDaJ/M++CRjur8vtXxoJnPETM8WRIc3CO0LyM=",
		Secret:         "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhdXRoLnRtcGltLnB3OnNlcnZlcjprZXk6ZmUyNzhhMmIwODY1NWM5YjYzNjQzZjI0N2E3NGFkMDciLCJpc3MiOiJhdXRoLnRtcGltLnB3OmNlbnRyYWwiLCJzZWNyZXQiOiJ6d2Q5TUpDVy9CbWdBcjNjeE0wbE1CU2tkaVVUa0JhUmNXZFp3ZkJPeGZRPSIsImlhdCI6MTY3MjczNTAxMSwic3ViIjoiZmUyNzhhMmIwODY1NWM5YjYzNjQzZjI0N2E3NGFkMDcifQ.RIpQebD1IgC7m7vtLzp_dDxNN0y6WSpU68PxlNBL9Ru7eFiP7hAbUwxX8X7B0PJv1MuRbJxrXpa2cwVChwIZfA",
		MiniServerHost: "http://localhost:4600",
		Debug:          true,
	}, tmpauth.FromHTTPHandler(mux))
	if err != nil {
		panic(err)
	}

	log.Println("listening on :4601")
	http.ListenAndServeTLS(":4601", "./192.168.17.145.pem", "./192.168.17.145-key.pem", ta.Stdlib())
}
