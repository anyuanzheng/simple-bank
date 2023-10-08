package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/iamzay/simplebank/db/sqlc"
	"github.com/iamzay/simplebank/util"
	"github.com/lib/pq"
)

type CreateUserArgs struct {
	Username string `json:"username" binding:"required,alphanum"`	
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`	
}

func getUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}	
}

func (server *Server) createUser(ctx *gin.Context) {
	// check args
	req	:= CreateUserArgs{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.GetHashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	createUserParams := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	}
	// call db
	user, err := server.store.CreateUser(ctx, createUserParams)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))	
		return
	}
	ctx.JSON(http.StatusOK, getUserResponse(user))
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshToken string `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User UserResponse
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	// check password
	if err := util.CheckPassword(req.Password, []byte(user.HashedPassword)); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// make token
	token, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return	
	}

	// make refresh token and insert to sessions table
	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return	
	}
	_, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID: refreshTokenPayload.ID,
		Username: user.Username,
		RefreshToken: refreshToken,
		UserAgent: ctx.Request.UserAgent(),
		ClientIp: ctx.ClientIP(),
		IsBlocked: false,
		ExpiresAt: refreshTokenPayload.ExpiredAt,	
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return	
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		AccessToken: token,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		User: getUserResponse(user),
	})
}
