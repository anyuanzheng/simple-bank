package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/iamzay/simplebank/db/sqlc"
	"github.com/iamzay/simplebank/token"
	"github.com/iamzay/simplebank/util"
)

type Server struct {
	store db.Store
	router *gin.Engine
	tokenMaker token.Maker
	config util.Config
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	server := &Server{ store: store, config: config, tokenMaker: tokenMaker}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.SetupRouter()
	
	return server, nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/login", server.loginUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfers", server.createTransfer)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{ "error": err.Error() }
}
