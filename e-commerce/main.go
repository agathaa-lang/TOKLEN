package main

import (
	"ecommerce/database"
	"ecommerce/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// koneksi database
	database.Connect()
	database.DB.AutoMigrate(&models.User{})

	r := gin.Default()

	// CORS (biar frontend connect)
	r.Use(cors.Default())

	// TEST
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server jalan",
		})
	})

	// =========================
	// REGISTER
	// =========================
	r.POST("/register", func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, gin.H{"message": "Data tidak valid"})
			return
		}

		// hash password
		hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		user.Password = string(hashed)

		// simpan ke DB
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(400, gin.H{"message": "Email sudah digunakan"})
			return
		}

		c.JSON(200, gin.H{"message": "Register berhasil"})
	})

	// =========================
	// LOGIN
	// =========================
	r.POST("/login", func(c *gin.Context) {
		var input models.User
		var user models.User

		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, gin.H{"message": "Data tidak valid"})
			return
		}

		// cari user
		if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
			c.JSON(400, gin.H{"message": "Email tidak ditemukan"})
			return
		}

		// cek password
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			c.JSON(400, gin.H{"message": "Password salah"})
			return
		}

		c.JSON(200, gin.H{"message": "Login berhasil"})
	})

	r.Run(":8080")
}
