package token

import (
	"testing"
	"time"

	"github.com/Ma-hiru/simplebank/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	var maker, err = NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	var username = util.RandomOwner()
	var duration = time.Minute
	var issuedAt = time.Now()
	var expiredAt = issuedAt.Add(duration)

	var token, err1 = maker.CreateToken(username, duration)
	require.NoError(t, err1)
	require.NotEmpty(t, token)

	var payload, err2 = maker.VerifyToken(token)
	require.NoError(t, err2)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	var maker, err = NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	var token, err1 = maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err1)
	require.NotEmpty(t, token)

	var payload, err2 = maker.VerifyToken(token)
	require.Error(t, err2)
	require.EqualError(t, err2, jwt.ErrTokenExpired.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlg(t *testing.T) {
	var payload, err = NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	var jwtToken = jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	var token, err1 = jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err1)

	var maker, err2 = NewJWTMaker(util.RandomString(32))
	require.NoError(t, err2)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, jwt.ErrTokenSignatureInvalid.Error())
	require.Nil(t, payload)
}
