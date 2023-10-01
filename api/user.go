package api

import (
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

type CreateUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`	
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
	
	// call db
	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	})
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
	ctx.JSON(http.StatusOK, CreateUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	})
}
