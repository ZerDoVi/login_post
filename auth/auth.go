package auth

import "golang.org/x/crypto/bcrypt"

func Hash(pass string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hashed)
}
func HashCheck(hashedpassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
	return err == nil
}
