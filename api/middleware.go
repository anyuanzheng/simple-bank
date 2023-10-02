package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iamzay/simplebank/token"
)

const (
	AuthorizationHeaderKey = "authorization"
	AuthorizationBearType = "Bearer"
	AuthorizationTokenCtxKey = "authorization_token"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	abort := func (ctx *gin.Context,  err error)  {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		ctx.Abort()
	}

	return func(ctx *gin.Context) {
		// get authorization header
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			abort(ctx, errors.New("authorization header is not provided"))	
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			abort(ctx, errors.New("authorization header format wrong"))	
			return	
		}

		// ensure it is bear type
		if fields[0] != AuthorizationBearType {
			abort(ctx, errors.New("authorization type is not supported"))	
			return
		}

		// verify token
		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			abort(ctx, errors.New("authorization invalid token"))	
			return	
		}

		// set token to ctx
		ctx.Set(AuthorizationTokenCtxKey, payload)
		ctx.Next()
	}
}
