package token

import (
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Payload contains the payload data of the token.
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Issuer    string    `json:"issuer"`
}

// GetExpirationTime 返回 token 的过期时间（exp），用于判断 token 是否已过期。
func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.ExpiredAt), nil
}

// GetIssuedAt 返回 token 的签发时间（iat），表示 token 何时被创建。
func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.IssuedAt), nil
}

// GetNotBefore 返回 token 的生效时间（nbf），表示 token 何时开始生效。
func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.IssuedAt), nil
}

// GetIssuer 返回 token 的签发者（iss），通常是签发 token 的服务标识。
func (payload *Payload) GetIssuer() (string, error) {
	return payload.Issuer, nil
}

// GetSubject 返回 token 的主题（sub），通常是 token 代表的用户或实体。
func (payload *Payload) GetSubject() (string, error) {
	return payload.Username, nil
}

// GetAudience 返回 token 的受众（aud），表示 token 预期的接收方。
func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}

func (payload *Payload) ToPasetoToken() paseto.Token {
	var token = paseto.NewToken()
	token.SetJti(payload.ID.String())
	token.SetSubject(payload.Username)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)
	token.SetNotBefore(payload.IssuedAt)
	token.SetIssuer(payload.Issuer)

	return token
}

// NewPayload creates a new token payload with a specific username and duration.
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	var tokenID, err = uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	var payload = &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
		Issuer:    "simplebank",
	}

	return payload, nil
}

func NewPayloadFromPasetoToken(token *paseto.Token) (*Payload, error) {
	ID, err := token.GetJti()
	if err != nil {
		return nil, err
	}
	username, err := token.GetSubject()
	if err != nil {
		return nil, err
	}
	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		return nil, err
	}
	expiredAt, err := token.GetExpiration()
	if err != nil {
		return nil, err
	}
	issuer, err := token.GetIssuer()
	if err != nil {
		return nil, err
	}
	tokenID, err := uuid.Parse(ID)
	if err != nil {
		return nil, err
	}

	var payload = &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
		Issuer:    issuer,
	}

	return payload, nil
}
