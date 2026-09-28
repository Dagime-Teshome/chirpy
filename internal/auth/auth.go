package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	is_pass, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return is_pass, err
	}
	return is_pass, nil
}
