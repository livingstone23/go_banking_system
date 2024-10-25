package db

import (
	"context"
	"go_banking_system/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	//hashes the password
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)


	arg := CreateUserParams{
		Username: util.RandomOwner(),
		HashedPassword:  hashedPassword,
		FullName: util.RandomOwner(),
		Email : util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	//require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	// Convierte pgtype.Timestamptz a time.Time
	createdAt1 := user1.CreatedAt.Time
	createdAt2 := user2.CreatedAt.Time

	require.WithinDuration(t, createdAt1, createdAt2, time.Second)
}