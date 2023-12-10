package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const jwtSecret string = "1234567890"

func TestUserCreateToken(t *testing.T) {
	var email = "email"
	var userToken = UserTokenService{}

	var token, err = userToken.Create(email, jwtSecret)

	if token == nil {
		t.Fatalf("returned token is null")
	}

	if err != nil {
		t.Fatalf("an error happened during token creation: %s", err)
	}
}

func TestUserValidateToken(t *testing.T) {
	var userToken = UserTokenService{}
	var email = "email"

	token, err := userToken.Create(email, jwtSecret)

	assert.Nil(t, err)

	decodedToken, err := userToken.Decode(*token, jwtSecret)
	assert.True(t, userToken.Validate(decodedToken))
}
