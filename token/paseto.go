package token

import (
	"time"

	"aidanwoods.dev/go-paseto"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	secretKey paseto.V4SymmetricKey
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey paseto.V4SymmetricKey) (*PasetoMaker, error) {
	return &PasetoMaker{secretKey: symmetricKey}, nil
}

// CreateToken creates a new token for a specific username and duration.
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	var payload, err = NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return payload.ToPasetoToken().V4Encrypt(maker.secretKey, nil), nil
}

// VerifyToken checks if the token is valid or not.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var parser = paseto.NewParser()
	var parseToken, err = parser.ParseV4Local(maker.secretKey, token, nil)
	if err != nil {
		return nil, err
	}

	return NewPayloadFromPasetoToken(parseToken)
}
