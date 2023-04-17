package microtoken

import (
	"encoding/base64"
	"log"
	"testing"
)

func TestToken(t *testing.T) {
	// This test data comes from a test application that is unused, and has its private key
	// revoked and the new one thrown away.
	var token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2ODE3MTQyNjMsInRva2VuIjoiZXlKaGJHY2lPaUpGVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SnBZWFFpT2pFMk9ERTNNVFF5TmpNc0ltbHpjeUk2SW1GMWRHZ3VkRzF3YVcwdWNIYzZZMlZ1ZEhKaGJDSXNJbUYxWkNJNkltRjFkR2d1ZEcxd2FXMHVjSGM2YzJWeWRtVnlPbWxrWlc1MGFYUjVPbVpsTWpjNFlUSmlNRGcyTlRWak9XSTJNelkwTTJZeU5EZGhOelJoWkRBM0lpd2ljM1ZpSWpvaVlXWTNaall5T1RZdFlqQTJOaTAwT1RrM0xXSmhNV0V0T1RCa01qY3pPREE0TVdFd0lpd2ljM1JoZEdWSlJDSTZJaUlzSW10cFpDSTZJbVZ6SW4wLkVib1pkUmJfVktPdnNzb0Zodll0MlZzX2VSVHFvaVBHRmZIOEdrc2puU293WXpwYllEOVBObjhYaUI2WnJWVFVNbjAtYm45WExQckpncWVScjhRZVZBIiwiaXNzIjoiYXV0aC50bXBpbS5wdzpjZW50cmFsOmZlMjc4YTJiMDg2NTVjOWI2MzY0M2YyNDdhNzRhZDA3IiwiYXVkIjoiYXV0aC50bXBpbS5wdzpzZXJ2ZXI6dXNlcl9jb29raWU6ZmUyNzhhMmIwODY1NWM5YjYzNjQzZjI0N2E3NGFkMDciLCJraWQiOiJlcyJ9.g9RkTgeRP5xryjaoMzhju2Iy6HAgitr18pDEfH9U1eo`

	codec := Codec{
		ClientID:   "fe278a2b08655c9b63643f247a74ad07",
		AuthDomain: "auth.tmpim.pw",
	}

	result, err := codec.EncodeToken([]byte(token))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	log.Println(base64.RawURLEncoding.EncodeToString(result))

	jwt, err := codec.DecodeToken(HS256Header, result)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	log.Println(string(jwt))

	if string(jwt) != token {
		t.Fatalf("token mismatch")
	}
}
