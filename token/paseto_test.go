package token

import (
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/Ma-hiru/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	var maker, err = NewPasetoMaker(paseto.NewV4SymmetricKey())
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

func TestExpiredPasetoToken(t *testing.T) {
	var maker, err = NewPasetoMaker(paseto.NewV4SymmetricKey())
	require.NoError(t, err)

	var token, err1 = maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err1)
	require.NotEmpty(t, token)

	var payload, err2 = maker.VerifyToken(token)
	require.Error(t, err2)
	require.EqualError(t, err2, "this token has expired")
	require.Nil(t, payload)
}
