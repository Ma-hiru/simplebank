package db

import (
	"context"
	"testing"
	"time"

	"github.com/Ma-hiru/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	var hashedPassword, err = util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	var arg = CreateUserParams{
		Username:     util.RandomOwner(),
		FullName:     util.RandomOwner(),
		HashPassword: hashedPassword,
		Email:        util.RandomEmail(),
	}
	{
		var user, err = testQueries.CreateUser(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, user)

		require.Equal(t, arg.Username, user.Username)
		require.Equal(t, arg.FullName, user.FullName)
		require.Equal(t, arg.HashPassword, user.HashPassword)
		require.Equal(t, arg.Email, user.Email)
		require.NotZero(t, user.PasswordChangedAt.IsZero())
		require.NotZero(t, user.CreatedAt)
		return user
	}
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	var user1 = createRandomUser(t)
	var user2, err = testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashPassword, user2.HashPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
