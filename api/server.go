package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kidx45/Project-KK/db/sqlc"
	"github.com/kidx45/Project-KK/token"
	"github.com/kidx45/Project-KK/utils"
)

type Server struct {
	Store  db.Store
	Router *gin.Engine
	Config utils.Config
	Token  token.Maker
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.SYMMETRIC_SECRET_KEY)
	if err != nil {
		panic("cannot create token maker: %w")
	}

	server := &Server{Store: store, Config: config, Token: tokenMaker}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", utils.ValidCurrency)
	}
	router.POST("/user", server.CreateUser)
	router.POST("/user/login", server.LoginUser)
	router.POST("refresh", server.Renew)

	authroutes := router.Group("/").Use(AuthMiddleware(server.Token))
	router.GET("/users", server.ListUsers)
	authroutes.GET("/user/:username", server.GetUserByUsername)

	authroutes.POST("/account", server.CreateAccount)
	authroutes.GET("/account/:id", server.GetAccount)
	authroutes.POST("/transfer", server.CreateTransfer)

	server.Router = router
	return server, nil
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

func (s *Server) Start(address string) error {
	return s.Router.Run(address)
}
