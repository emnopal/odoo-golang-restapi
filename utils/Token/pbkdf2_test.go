package utils

import (
	"testing"
)

func TestPbkdf2Encoder(t *testing.T) {
	password := "admin123"
	hash, err := Pbkdf2Encoder(password)
	if err != nil {
		t.Error(err)
	}
	if hash == "" {
		t.Error(`hash expected: $pbkdf2-sha512$25000$<salt>$<hash>, found ""`)
	}
	t.Log(hash)
}

func TestPbkdf2DecoderValid(t *testing.T) {
	password := "admin123"
	hash := "$pbkdf2-sha512$25000$xZgT4jxnTGlNKSVESEnpXQ$3Y18V7rxYPfRzl6EnTJba1ITGGMKdjqH.xl.WpTj4KGqCR3g/NRPyWqnB1sd6UdYdnTOPuxZ3t5Ilwi7FNodVg"
	isValid, err := Pbkdf2Decoder(password, hash)
	if err != nil {
		t.Error(err)
	}
	if !isValid {
		t.Errorf(`hash expected valid with %s, found false`, password)
	}
	t.Log(isValid)
}

func TestPbkdf2DecoderInvalid(t *testing.T) {
	password := "admin123"
	hash := "$pbkdf2-sha512$25000$xZgT4jxnTGlNKSVESEnpxQ$3Y18V7rxYPfRzl6EnTJba1ITGGMKdjqH.xl.WpTj4KGqCR3g/NRPyWqnB1sd6UdYdnTOPuxZ3t5Ilwi7FNodVg"
	isValid, err := Pbkdf2Decoder(password, hash)
	if err == nil {
		t.Errorf(`hash expected error since it's invalid with %s, found not error`, password)
	}
	if isValid {
		t.Errorf(`hash expected invalid with %s, found true`, password)
	}
	t.Log(isValid)
}

func TestPbkdf2EncDec01(t *testing.T) {
	password := "admin123"
	hash, err := Pbkdf2Encoder(password)
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
	isValid, err := Pbkdf2Decoder(password, hash)
	if err != nil {
		t.Error(err)
	}
	t.Log(isValid)
}

func TestPbkdf2EncDec02(t *testing.T) {
	password := "123456"
	hash, err := Pbkdf2Encoder(password)
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
	isValid, err := Pbkdf2Decoder(password, hash)
	if err != nil {
		t.Error(err)
	}
	t.Log(isValid)
}

func TestPbkdf2EncDec03(t *testing.T) {
	password := "asd4567890"
	hash, err := Pbkdf2Encoder(password)
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
	isValid, err := Pbkdf2Decoder(password, hash)
	if err != nil {
		t.Error(err)
	}
	t.Log(isValid)
}
