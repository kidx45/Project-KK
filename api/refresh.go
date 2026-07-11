package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RenewRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiredAt time.Time `json:"access_token_expired_at"`
}

func (s *Server) Renew(ctx *gin.Context) {
	var req RenewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshTokenPayload, err := s.Token.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := s.Store.GetSessionById(ctx, refreshTokenPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if refreshTokenPayload.Username != session.Username {
		err := fmt.Errorf("incorrect username in token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if req.RefreshToken != session.RefreshToken {
		err := fmt.Errorf("mismatched refresh token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.ExpiresAt.Before(time.Now()) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessTokenPayload, err := s.Token.CreateToken(
		refreshTokenPayload.Username,
		s.Config.ACCESS_TOKEN_DURATION,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rps := RenewResponse{
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  accessTokenPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rps)
}