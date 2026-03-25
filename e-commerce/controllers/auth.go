package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User

	// ambil JSON dari frontend
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Data tidak valid",
		})
		return
	}

	// hash password
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hash)

	// simpan ke DB
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email sudah digunakan",
		})
		return
	}

	// sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Register berhasil",
	})
}
