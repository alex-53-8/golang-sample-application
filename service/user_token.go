package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserTokenDecoded struct {
	Email   string
	Expires time.Time
}

type UserToken interface {
	Create(email string, secret string) (*string, error)
	Decode(token string, secret string) (*UserTokenDecoded, error)
	Validate(decodedToken *UserTokenDecoded) bool
}

type UserTokenService struct{}

func (u UserTokenService) Create(email string, secret string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   email,
		"expires": time.Now().Add(time.Hour * 24 * 10).Format(time.RFC3339Nano),
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (u UserTokenService) Validate(decodedToken *UserTokenDecoded) bool {
	if decodedToken == nil {
		return false
	}

	return time.Now().Before(decodedToken.Expires)
}

func (u UserTokenService) Decode(tokenString string, secret string) (*UserTokenDecoded, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(secret), nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var userToken = UserTokenDecoded{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email := fmt.Sprintf("%v", claims["email"])

		time, err := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", claims["expires"]))
		if err != nil {
			log.Println(err)
			return nil, errors.New("expiration time from token is not recognized")
		}

		userToken.Email = email
		userToken.Expires = time
	} else {
		return nil, err
	}

	return &userToken, nil
}
