package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/iamzay/simplebank/db/sqlc"
	"github.com/iamzay/simplebank/util"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey: util.RandomString(32),
		TokenDuration: 15*time.Minute,
	}
	server, err := NewServer(store, config)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
