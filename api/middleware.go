package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kidx45/Project-KK/Backend-Team/token"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("authorization")
		if authorizationHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("no authorization header")))
			return
		}

		authorizationSplit := strings.Split(authorizationHeader, " ")
		if len(authorizationSplit) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid authorization header format")))
			return
		}

		if authorizationSplit[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid authorization header format")))
			return
		}

		token := authorizationSplit[1]
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("no token")))
			return
		}

		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set("auth_payload", payload)
		ctx.Next()
	}
}
