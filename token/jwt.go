package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration.
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	var payload, err = NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	var jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	var jwtToken, err = jwt.ParseWithClaims(token, &Payload{}, maker.keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return nil, err
	}

	var payload, ok = jwtToken.Claims.(*Payload)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return payload, nil
}

// keyFunc checks the signing method and payload ,then returns the secret key for validation.
func (maker *JWTMaker) keyFunc(token *jwt.Token) (any, error) {
	var _, ok = token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, jwt.ErrTokenSignatureInvalid
	}

	var exp, err = token.Claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}
	if time.Now().After(exp.Time) {
		return nil, jwt.ErrTokenExpired
	}

	return []byte(maker.secretKey), nil
}
