package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/kidx45/Project-KK/Backend-Team/token"
)

type CreateTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) CreateTransfer(ctx *gin.Context) {
	var req CreateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, err := s.Store.GetAccountById(ctx, req.FromAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := ctx.MustGet("auth_payload").(*token.Payload)
	if payload.Username != fromAccount.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("From Account doesn't belong to user")))
		return
	}
	toAccount, err := s.Store.GetAccountById(ctx, req.ToAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if toAccount.Currency != req.Currency || fromAccount.Currency != req.Currency {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("Currency doesn't match")))
		return
	}

	result, err := s.Store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}
