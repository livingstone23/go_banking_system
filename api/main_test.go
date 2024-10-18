package api

import (

	"os"
	"testing"
	"github.com/gin-gonic/gin"
)


func TestMain(m *testing.M) {
	
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
    
    // Exit with the test code
    os.Exit(m.Run())

}