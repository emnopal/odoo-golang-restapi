package testing

import (
	"testing"

	jwt_schema "github.com/emnopal/odoo-golang-restapi/app/schemas/auth"
	Token "github.com/emnopal/odoo-golang-restapi/app/utils/Token"
)

func TestGetBearerToken(t *testing.T) {
	Header := "Authorization eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2OTA5NjM4NTYsImlkIjoiMSIsInVzZXJuYW1lIjoiVXNlciBUZXN0In0.S9xh83Xzswyk1nmddGOSnZlMbbopnqsW3LGnYRuS0Us"
	bearerToken := Token.GetBearerToken(Header)
	if bearerToken == "" {
		t.Error(`Expected bearer token, found ""`)
	}
	t.Log(bearerToken)
}

func TestGetBearerTokenInvalid01(t *testing.T) {
	Header := "Authorization"
	bearerToken := Token.GetBearerToken(Header)
	if bearerToken != "" {
		t.Errorf(`Expected "", found %s`, bearerToken)
	}
}

func TestGetBearerTokenInvalid02(t *testing.T) {
	Header := "Authorization "
	bearerToken := Token.GetBearerToken(Header)
	if bearerToken != "" {
		t.Errorf(`Expected "", found %s`, bearerToken)
	}
}

func TestGenerateToken(t *testing.T) {
	jwt_claims := jwt_schema.JWTAccessClaims{
		ID:       "1",
		Username: "User Test",
	}
	token, err := Token.Generate(&jwt_claims)
	if err != nil {
		t.Errorf("Expected jwt token, error found. Error message: %s", err)
	}
	t.Log(token)
}

func TestParseBearerToken(t *testing.T) {
	jwt_claims := jwt_schema.JWTAccessClaims{
		ID:       "1",
		Username: "User Test",
	}
	token, err := Token.Generate(&jwt_claims)
	if err != nil {
		t.Errorf("Expected jwt token, error found. Error message: %s", err)
	}
	bearerToken, err := Token.ParseBearerToken(token)
	if err != nil {
		t.Errorf("Expected jwt token, error found. Error message: %s", err)
	}
	t.Log(bearerToken)
}

func TestExtractID(t *testing.T) {
	jwt_claims := jwt_schema.JWTAccessClaims{
		ID:       "1",
		Username: "User Test",
	}
	token, err := Token.Generate(&jwt_claims)
	if err != nil {
		t.Errorf("Expected jwt token, error found. Error message: %s", err)
	}
	id, err := Token.ExtractID(token)
	if err != nil {
		t.Errorf("Expected jwt token, error found. Error message: %s", err)
	}
	if id != jwt_claims.ID {
		t.Error("id is not same")
	}
	t.Log(id)
	t.Log(jwt_claims.ID)
}

func TestExtractUserName(t *testing.T) {
	jwt_claims := jwt_schema.JWTAccessClaims{
		ID:       "1",
		Username: "User Test",
	}
	token, err := Token.Generate(&jwt_claims)
	if err != nil {
		t.Errorf("Expected jwt token, error found. Error message: %s", err)
	}
	uname, err := Token.ExtractUserName(token)
	if err != nil {
		t.Errorf("Expected jwt token, error found. Error message: %s", err)
	}
	if uname != jwt_claims.Username {
		t.Error("username is not same")
	}
	t.Log(uname)
	t.Log(jwt_claims.Username)
}

func TestExtractJWTClaims(t *testing.T) {
	jwt_claims := &jwt_schema.JWTAccessClaims{
		ID:       "1",
		Username: "User Test",
	}
	token, err := Token.Generate(jwt_claims)
	if err != nil {
		t.Errorf("Expected jwt token, error found. Error message: %s", err)
	}
	claims, err := Token.ExtractJWTClaims(token)
	if err != nil {
		t.Errorf("Expected jwt token, error found. Error message: %s", err)
	}
	if (jwt_claims.ID != claims.ID) && (jwt_claims.Username != claims.Username) {
		t.Error("JWT Claims is not same")
	}
	t.Log(jwt_claims)
	t.Log(claims)
}

func TestValidToken(t *testing.T) {
	jwt_claims := &jwt_schema.JWTAccessClaims{
		ID:       "1",
		Username: "User Test",
	}
	token, err := Token.Generate(jwt_claims)
	if err != nil {
		t.Errorf("Expected jwt token, error found. Error message: %s", err)
	}
	err = Token.Validate(token)
	if err != nil {
		t.Errorf("Expected valid jwt token, error found. Error message: %s", err)
	}
}

func TestInvalidToken(t *testing.T) {
	jwt_claims := &jwt_schema.JWTAccessClaims{
		ID:       "1",
		Username: "User Test",
	}
	token, err := Token.Generate(jwt_claims)
	token = token + "."
	if err != nil {
		t.Errorf("Expected jwt token, error found. Error message: %s", err)
	}
	err = Token.Validate(token)
	if err == nil {
		t.Error("Expected invalid, found not invalid")
	}
}

func TestInvalidToken02(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	err := Token.Validate(token)
	if err == nil {
		t.Error("Expected invalid, found not invalid")
	}
}
