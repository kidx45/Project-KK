package test

import (
	"testing"

	"github.com/kidx45/Project-KK/Backend-Team/utils"
	"github.com/stretchr/testify/require"
)

func TestPassword (t *testing.T) {
	password := utils.RandomLetterString(10)
	hashed,err := utils.HashPassword(password)
	require.Nil(t,err)
	require.NotEmpty(t,hashed)
	require.NotEqual(t,password,hashed)

	err = utils.CheckPassword(password,hashed)
	require.NoError(t,err)
}