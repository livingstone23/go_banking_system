package api

import (
	db "go_banking_system/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for banking services.
type Server struct {
	store  *db.Store   //it will allow us to interact with the database
	router *gin.Engine //it will allow routing of HTTP requests to the appropriate handler functions
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// add routes to router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Create general function to handle the errors response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
