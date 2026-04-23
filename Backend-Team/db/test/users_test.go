package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	sqlc "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	utils "github.com/kidx45/Project-KK/Backend-Team/utils"
)

func createRandomUser(t *testing.T) sqlc.User {
	arg := sqlc.CreateUserParams{
		Username: utils.RandomUserName(),
		HashedPassword: utils.RandomPassword(),
		Email: utils.RandomEmail(),
		FullName: utils.RandomFullName(),
		PhoneNumber: utils.RandomPhoneNumber(),
	}
	
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
	return user
}
		
func TestCreateUser(t *testing.T) {
	createRandomUser(t) 
}
