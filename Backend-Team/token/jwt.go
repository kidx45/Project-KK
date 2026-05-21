package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTToken (secretKey string) (Maker,error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("Secret Key size must be %v",minSecretKeySize)
	}
	return &JWTMaker{secretKey: secretKey},nil
}

func (j *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username,duration)
	if err != nil {
		return "", fmt.Errorf("cant't create token %s",err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,payload)
	return token.SignedString([]byte(j.secretKey))
}
func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{},error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid Token")
		}
		return []byte(j.secretKey),nil
	}

	jwtToken, err := jwt.ParseWithClaims(token,&Payload{},keyFunc)
	if err != nil {
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr, errors.New("Expired Token")) {
			return nil, errors.New("Expired Token")
		}
		return nil, errors.New("Invalid Token")
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errors.New("Invalid Token")
	}
	return payload, nil
}