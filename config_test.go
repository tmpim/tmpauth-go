package tmpauth

import (
	"encoding/json"
	"testing"
)

func TestConfig(t *testing.T) {
	var content = `{
	"publicKey": "BN/PHEYgs0meH878gqpWl81WD3zEJ+ubih3RVYwFxaYXxHF+5tgDaJ/M++CRjur8vtXxoJnPETM8WRIc3CO0LyM=",
    "secret": "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhdXRoLnRtcGltLnB3OnNlcnZlcjprZXk6NDE3ZGE5NmM0YTE0ODY1YjhkZDg2NTczYjkwODk0NzYiLCJpc3MiOiJhdXRoLnRtcGltLnB3OmNlbnRyYWwiLCJzZWNyZXQiOiJqR1dwaHpaYTQ3Y0pJMEYvZzh5dEFXUTBwWnQwQlhwZkFtYzRXL0VSRUtrPSIsImlhdCI6MTY3MTQxODUwOSwic3ViIjoiNDE3ZGE5NmM0YTE0ODY1YjhkZDg2NTczYjkwODk0NzYifQ.0scVeVwKbPO-CjcSBB7M1gMKLjHmVt3AdjbtaDpog5oC6GDB78GfcRpGHflW5CmHAup4TQHhKpm3h_trc4TlMw"
	}`

	var config UnserializableConfig
	err := json.Unmarshal([]byte(content), &config)
	if err != nil {
		t.Fatal(err)
	}

	_, err = config.Parse()
	if err != nil {
		t.Fatal(err)
	}
}
