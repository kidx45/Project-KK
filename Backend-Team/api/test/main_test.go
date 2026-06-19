package test

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kidx45/Project-KK/Backend-Team/api"
	"github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/kidx45/Project-KK/Backend-Team/utils"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *api.Server {
	config := utils.Config{
		SYMMETRIC_SECRET_KEY:  utils.RandomLetterString(32),
		ACCESS_TOKEN_DURATION: time.Minute,
	}
	server, err := api.NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
