package sessions

import (
	"crypto/rand"
	"encoding/hex"
)

func Session() string {
	a := make([]byte, 32)
	rand.Read(a)
	return hex.EncodeToString(a)
}
