package test

import (
	"testing"
	"time"

	"github.com/kidx45/Project-KK/Backend-Team/token"
	"github.com/kidx45/Project-KK/Backend-Team/utils"
	"github.com/stretchr/testify/require"
)

func TestJWTToken(t *testing.T) {
	maker, err := token.NewJWTToken(utils.RandomLetterString(32))
	require.NoError(t,err)

	username := utils.RandomLetterString(10)
	duration := time.Minute

	token,err := maker.CreateToken(username,duration)
	require.NoError(t,err)
	require.NotEmpty(t,token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t,err)
	require.NotEmpty(t,payload)

	require.Equal(t,username,payload.Username)
	require.WithinDuration(t,time.Now(),payload.ExpiredAt,duration)
}

func TestMinSecretKeySize(t *testing.T) {
	maker,err := token.NewJWTToken(utils.RandomLetterString(31))
	require.Error(t,err)
	require.Nil(t,maker)
}

func TestExpiredToken(t *testing.T) {
	maker,err := token.NewJWTToken(utils.RandomLetterString(32))
	require.NoError(t,err)

	token,err := maker.CreateToken("username",-time.Minute)
	require.NoError(t,err)
	require.NotEmpty(t,token)

	payload, err := maker.VerifyToken(token)
	require.Error(t,err)
	require.Nil(t,payload)
}