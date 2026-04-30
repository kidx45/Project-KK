package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/kidx45/Project-KK/Backend-Team/utils"
)

type Server struct {
	Store  db.Store
	Router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{Store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency",utils.ValidCurrency)
	}
	router.POST("/user",server.CreateUser)
	router.GET("/users",server.ListUsers)
	router.GET("/user/:username",server.GetUserByUsername)
	router.POST("/transfer",server.CreateTransfer)

	server.Router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

func (s *Server) Start (address string) error {
	return s.Router.Run(address)
}
