package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewTokenResponse struct {
	AccessToken string `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewToken(ctx *gin.Context) {
	var request renewTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return		
	}

	// verify refresh token
	payload, err := server.tokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// get session by refresh token id
	session, err := server.store.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return	
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != payload.Username || session.RefreshToken != request.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("mismatch token")))
		return	
	}
	
	// create a new access token and return
	token, accessTokenPayload, err := server.tokenMaker.CreateToken(payload.Username, server.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return	
	}

	ctx.JSON(http.StatusOK, renewTokenResponse{
		AccessToken: token,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	})
}
