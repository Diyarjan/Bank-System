package handler

import (
	"github.com/Diyarjan/BankSystem"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) createAccount(c *gin.Context) {
	var account BankSystem.ToMakeAccount

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid Input structure ": err.Error()})
		return
	}

	accountID, err := h.services.Control.CreateAccount(account)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error from create account": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"account created your id": accountID})
}

func (h *Handler) deleteAccount(c *gin.Context) {

	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid account_id": err.Error()})
	}

	if err = h.services.Control.DeleteAccount(accountID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error from delete account": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"deleted account id - ": accountID})

}

func (h *Handler) getAccountByID(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid account_id": err.Error()})
	}

	account, err := h.services.GetAccountByID(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error from get account by id": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})

}

func (h *Handler) getAllAccounts(c *gin.Context) {
	var accounts []BankSystem.Account
	accounts, err := h.services.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error from get all accounts": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)

}
