package jwt_test

import (
	"adv-go/api/pkg/jwt"
	"testing"
)

const (
	secret = "a1df53b4837446a7254dcdb4463848f8"
	email = "a@a1.ru"
)

func TestJWTCreated(t *testing.T) {
	jwtService := jwt.NewJWT(secret)
	token, err := jwtService.Create(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
	}
	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("not valid token, token: ", token)
	}
	if data.Email != email {
		t.Fatalf("Wrong data. Email want: %s, got: %s", data.Email, email)
	}
}