package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
)

type CreateUserRequest struct {
	Username       string `json:"username" binding:"required"`
	HashedPassword string `json:"hashedPassword" binding:"required"`
	Email          string `json:"email" binding:"required"`
	FullName       string `json:"fullName" binding:"required"`
	PhoneNumber    string `json:"phoneNumber" binding:"required"`
}

func (s *Server) CreateUser (ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		HashedPassword: req.HashedPassword,
		Email: req.Email,
		FullName: req.FullName,
		PhoneNumber: req.PhoneNumber,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated,user)
}

type ListUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) ListUsers (ctx *gin.Context) {
	var req ListUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := s.store.ListUsers(ctx,db.ListUsersParams{
		Limit: req.PageSize,
		Offset: (req.PageID-1)*req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return 
	}

	ctx.JSON(http.StatusOK,users)
}