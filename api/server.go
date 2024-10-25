package api

import (
	"fmt"
	db "go_banking_system/db/sqlc"
	"go_banking_system/token"
	"go_banking_system/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for banking services.
type Server struct {
	config util.Config //it will store the configuration of the server
	store  db.Store   //it will allow us to interact with the database
	tokenMaker token.Maker
	router *gin.Engine //it will allow routing of HTTP requests to the appropriate handler functions
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server,error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}


	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	
	return server, nil
}

func (server *Server) setupRouter(){

	router := gin.Default()

	// add routes to router
	router.POST("/free_accounts", server.createAccount)
	router.GET("/free_accounts/:id", server.getAccount)
	router.GET("/free_accounts", server.listAccounts)


	// Routes for the transfer
	router.POST("/free_transfers", server.createTransfer)

	// Routes for the users
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// Routes for the auth middleware
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// Router with auth middleware
		// Routes for the accounts
		authRoutes.POST("/accounts", server.createAccount)
		authRoutes.GET("/accounts/:id", server.getAccount)
		authRoutes.GET("/accounts", server.listAccounts)

		// Routes for the transfer
		authRoutes.POST("/transfers", server.createTransfer)



	server.router = router

}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Create general function to handle the errors response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
