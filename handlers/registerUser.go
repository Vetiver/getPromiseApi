package handlers

import (
	"net/http"
	"getPromiseApi/db"
	"github.com/gin-gonic/gin"
)

func (h BaseHandler) RegisterUser(c *gin.Context) {
	var user *db.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	registeredUser, err := h.db.RegisterUser(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
        "result": registeredUser,
    })
}
