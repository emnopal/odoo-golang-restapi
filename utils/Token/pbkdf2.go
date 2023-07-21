package utils

import (
	"gopkg.in/hlandau/passlib.v1"
	"gopkg.in/hlandau/passlib.v1/abstract"
	"gopkg.in/hlandau/passlib.v1/hash/pbkdf2"
)

func Pbkdf2Encoder(password string) (hash string, err error) {
	var passlib passlib.Context
	passlib.Schemes = []abstract.Scheme{
		pbkdf2.SHA512Crypter,
	}
	hash, err = passlib.Hash(password)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func Pbkdf2Decoder(password, encryptedPassword string) (verify bool, err error) {
	var passlib passlib.Context
	passlib.Schemes = []abstract.Scheme{
		pbkdf2.SHA512Crypter,
	}
	ver, err := passlib.Verify(password, encryptedPassword)
	if ver == "" && err == nil {
		return true, nil
	}
	return false, err
}
