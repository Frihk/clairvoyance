package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a bcrypt hash of the provided plaintext password.
// Use this when creating or updating a user's password before saving to the DB.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a plaintext password with a bcrypt hashed password and
// returns true when they match. Use this during login to validate credentials.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
