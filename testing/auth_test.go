package testing

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashCheck(t *testing.T) {
	pass := "12345"
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to create hash")
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		t.Fatalf("Comparison failed")
	}

}

func TestWrongHashCheck(t *testing.T) {
	pass := "12345"
	wrongpass := "12346"
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("(wrong) Failed to create hash")
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(wrongpass))
	if err != nil {
		t.Fatalf("(wrong) Comparison failed")
	}

}
