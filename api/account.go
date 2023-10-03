package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/iamzay/simplebank/db/sqlc"
	"github.com/iamzay/simplebank/token"
	"github.com/lib/pq"
)

type CreateAccountArgs struct {
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	// check args
	req	:= CreateAccountArgs{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	token := ctx.MustGet(AuthorizationTokenCtxKey).(*token.Payload)
	// call db
	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{ Owner: token.Username, Currency: req.Currency, Balance: 0 })
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))	
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type GetAccountArgs struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}
func (server *Server) getAccount(ctx *gin.Context) {
	req := GetAccountArgs{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return	
	}
	account, err := server.store.GetAccount(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))	
		return
	}

	token := ctx.MustGet(AuthorizationTokenCtxKey).(*token.Payload)
	if account.Owner != token.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)	
}

type ListAccountRequest struct {
	PageId int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}
func (server *Server) listAccount(ctx *gin.Context) {
	var req ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	token := ctx.MustGet(AuthorizationTokenCtxKey).(*token.Payload)
	var listAccountParams = db.ListAccountsParams{
		Owner: token.Username,
		Limit: req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, listAccountParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
