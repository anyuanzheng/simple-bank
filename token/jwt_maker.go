package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMaker struct {
	secretKey string
}

const minSecretKeySize = 32

func NewJWTMaker(secretkey string) (Maker, error) {
	if len(secretkey) < minSecretKeySize {
		return nil, errors.New("secret key is too short")
	}	
	return &JWTMaker{ secretKey: secretkey }, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload,err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(maker.secretKey))
	return signedToken, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	parsedToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := parsedToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}