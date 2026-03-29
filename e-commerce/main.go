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

	// CORS
	r.Use(cors.Default())

	// HTML & STATIC
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// =========================
	// HALAMAN WEB (FIX 404)
	// =========================
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})

	r.GET("/products", func(c *gin.Context) {
		c.HTML(200, "products.html", nil)
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

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			c.JSON(500, gin.H{"message": "Gagal hash password"})
			return
		}

		user.Password = string(hashed)

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

		if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
			c.JSON(400, gin.H{"message": "Email tidak ditemukan"})
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			c.JSON(400, gin.H{"message": "Password salah"})
			return
		}

		c.JSON(200, gin.H{"message": "Login berhasil"})
	})

	// RUN SERVER
	r.Run(":8080")
}