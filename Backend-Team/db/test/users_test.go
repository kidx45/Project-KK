package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	sqlc "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	
)

func TestCreateUser(t *testing.T) {
	arg := sqlc.CreateUserParams{
		Username: "kidx45",
		HashedPassword: "some_hashed_password",
		Email: "some_email@gmail.com",
		FullName: "some_full_name",
		PhoneNumber: "0909090909",
	}
	
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
}