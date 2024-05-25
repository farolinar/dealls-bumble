package password

import (
	"errors"

	"github.com/farolinar/dealls-bumble/config"
	"golang.org/x/crypto/bcrypt"
)

func Hash(plaintextPassword string) (string, error) {
	cfg := config.GetConfig()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), cfg.App.BCryptSalt)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func Matches(plaintextPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
