package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/lib/pq"
	"github.com/kidx45/Project-KK/Backend-Team/utils"
)

type CreateUserRequest struct {
	Username       string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	FullName       string `json:"fullName" binding:"required"`
	PhoneNumber    string `json:"phoneNumber" binding:"required"`
}

type CreateUserResponse struct {
	User string `json:"username"`
	Email string `json:"email"`
	FullName string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
}

func (s *Server) CreateUser (ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	user, err := s.Store.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashed,
		Email: req.Email,
		FullName: req.FullName,
		PhoneNumber: req.PhoneNumber,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated,CreateUserResponse{
		User: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		PhoneNumber: user.PhoneNumber,
	})
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

	users, err := s.Store.ListUsers(ctx,db.ListUsersParams{
		Limit: req.PageSize,
		Offset: (req.PageID-1)*req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return 
	}

	ctx.JSON(http.StatusOK,users)
}

type GetUserRequest struct {
	Username string `uri:"username" binding:"required"`
}

func (s *Server) GetUserByUsername (ctx *gin.Context) {
	var req GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	user,err := s.Store.GetUser(ctx,req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound,errorResponse(err))
		return
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,user)
}