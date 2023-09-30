package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/iamzay/simplebank/db/sqlc"
)

type CreateAccountArgs struct {
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	// check args
	req	:= CreateAccountArgs{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// call db
	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{ Owner: req.Owner, Currency: req.Currency, Balance: 0 })
	if err != nil {
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
	var listAccountParams = db.ListAccountsParams{
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
