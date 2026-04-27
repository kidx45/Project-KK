package api

import (
	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/user",server.CreateUser)
	router.GET("/users",server.ListUsers)

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

func (s *Server) Start (address string) error {
	return s.router.Run(address)
}
