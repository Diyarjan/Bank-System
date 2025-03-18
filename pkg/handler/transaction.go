package handler

import (
	"errors"
	"fmt"
	"github.com/Diyarjan/BankSystem"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Deposit struct {
	Amount float32 `json:"amount" binding:"required"`
}
type Withdraw struct {
	Amount float32 `json:"amount" binding:"required"`
}

func (h *Handler) depositToAccount(c *gin.Context) {
	var deposit Deposit
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid account_id": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&deposit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	depositStruct := BankSystem.DebitCreditStruct{
		AccountID: accountID,
		Amount:    deposit.Amount,
	}

	err = h.services.Transaction.DepositToAccount(depositStruct)
	if err != nil {
		if errors.Is(err, errors.New("not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("account %d not found", accountID)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("account %d deposit successfully", accountID))
}

func (h *Handler) withdrawFromAccount(c *gin.Context) {
	var withdraw Withdraw

	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid account_id": err.Error()})
	}

	if err := c.ShouldBindJSON(&withdraw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	withdrawStruct := BankSystem.DebitCreditStruct{
		AccountID: accountID,
		Amount:    withdraw.Amount,
	}

	if err := h.services.Transaction.WithdrawFromAccount(withdrawStruct); err != nil {
		if errors.Is(err, errors.New("not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("account %d not found", accountID)})
			return
		}
		if errors.Is(err, errors.New("not enough funds")) {
			c.JSON(http.StatusOK, gin.H{"error": fmt.Errorf("account %d not enough funds", accountID)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("account %d withdraw successfully", accountID))
}

func (h *Handler) transferFunds(c *gin.Context) { // transferToAccount
	var transfer BankSystem.Transfer

	fromAccountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid account_id": err.Error()})
		return
	}
	if err = c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	transfer.FromAccountId = fromAccountID

	balance, err := h.services.Transaction.TransferFunds(transfer)

	if err != nil {
		if errors.Is(err, errors.New("not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("some account not found")})
			return
		}
		if errors.Is(err, errors.New("not enough funds")) {
			c.JSON(http.StatusOK, gin.H{"error": fmt.Errorf("account %d not enough funds, your balance %d", fromAccountID, balance)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Transfer successfully your balance": balance})
}

func (h *Handler) getTransactionHistory(c *gin.Context) {
	var transactions []BankSystem.Transaction
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid account_id": err.Error()})
		return
	}

	transactions, err = h.services.Transaction.GetTransactionHistory(accountID)
	if err != nil {
		if errors.Is(err, errors.New("not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("account %d not found", accountID)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(transactions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "transactions not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
