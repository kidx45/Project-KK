package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < 32 {
		return nil, errors.New("symmetric key must be at least 32 bytes long")
	}
	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}
func (j *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", fmt.Errorf("cant't create token %s", err)
	}
	return j.paseto.Encrypt(j.symmetricKey, payload, nil)
}
func (j *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := j.paseto.Decrypt(token, j.symmetricKey, payload, nil)
	if err != nil {
		return nil, errors.New("Invalid Token")
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}
