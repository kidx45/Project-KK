package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/kidx45/Project-KK/Backend-Team/token"
)

type CreateAccountRequest struct {
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) CreateAccount (ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet("auth_payload").(*token.Payload)

	account, err := s.Store.CreateAccount(ctx,db.CreateAccountParams{
		Username: payload.Username,
		Balance: 100,
		Currency: req.Currency,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK,account)
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) GetAccount (ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account,err := s.Store.GetAccountById(ctx,req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}
	
	payload := ctx.MustGet("auth_payload").(*token.Payload)
	if payload.Username != account.Username {
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)
}

type ListAccountsRequest struct {
	PageID int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5"`
}

func (s *Server) ListAccounts (ctx *gin.Context) {
	var req ListAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}
	payload := ctx.MustGet("auth_payload").(*token.Payload)
	accounts,err := s.Store.ListAccounts(ctx,db.ListAccountsParams{
		Username: payload.Username,
		Limit: int32(req.PageSize),
		Offset: int32((req.PageID - 1) * req.PageSize),
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,accounts)
}