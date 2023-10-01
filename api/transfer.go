package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/iamzay/simplebank/db/sqlc"
)

type createTransferRequest struct {
	FromAccountId int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountId int64 `json:"to_account_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.isAccountValid(ctx, req.FromAccountId, req.Currency) {
		return
	}
	if !server.isAccountValid(ctx, req.ToAccountId, req.Currency) {
		return
	}

	result, err := server.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID: req.ToAccountId,
		Amount: req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) isAccountValid(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency not match: %s %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
