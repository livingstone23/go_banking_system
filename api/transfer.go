package api

import (
	"database/sql"
	"fmt"
	db "go_banking_system/db/sqlc"
	"go_banking_system/token"
	"net/http"

	"github.com/gin-gonic/gin"
)



type transferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

// createAccount
func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	

	//check if the account is valid
	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency) 
	if !valid {
		return
	}

	//with the authoriation middleware
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := fmt.Errorf("from account %d does not belong to the authenticated user", req.FromAccountID)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	// create a new account
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

//Func help to valid the account currency
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
			//check if not found register
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return account, false
			}

			//check if other error
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return account, false
	}

	//check if the currency is the same
	if account.Currency != currency {
		err := fmt.Errorf("account %d currency mismatch %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, err)
		return account, false
	}

	return account, true
	

}