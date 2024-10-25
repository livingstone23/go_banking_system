package api

import (

	"os"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "go_banking_system/db/sqlc"
	"go_banking_system/util"
	"time"




)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}


func TestMain(m *testing.M) {
	
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
    
    // Exit with the test code
    os.Exit(m.Run())

}

