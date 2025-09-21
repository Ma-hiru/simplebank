package api

import (
	db "github.com/Ma-hiru/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store db.Store) *Server {
	var server = &Server{store: store, router: gin.Default()}

	configureValidator()
	configureRouter(server)

	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func configureValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var err = v.RegisterValidation("currency", validCurrency)
		if err != nil {
			panic(err)
		}
	}
}

func configureRouter(server *Server) {
	server.router.POST("/accounts", server.createAccount)
	server.router.GET("/accounts", server.listAccount)
	server.router.PUT("/accounts", server.updateAccount)
	server.router.GET("/accounts/:id", server.getAccount)
	server.router.DELETE("/accounts/:id", server.deleteAccount)

	server.router.POST("/transfers", server.createTransfer)
}

func errResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
