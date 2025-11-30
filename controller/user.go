package controller

import (
	"API/config"
	"API/models"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserCacheKey = "all_users"
	UserCacheTTL = 5 * time.Minute
)

func GetUsers(c *gin.Context) {
	ctx := c.Request.Context()

	if config.RedisClient != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		cacheData, err := config.RedisClient.Get(cacheCtx, UserCacheKey).Result()
		// redis.Nil means cache miss, which is not a real error for us.
		if err == nil { // Cache hit
			user := []models.User{}
			if err := json.Unmarshal([]byte(cacheData), &user); err == nil {
				c.JSON(http.StatusOK, gin.H{"source": "cache", "data": user})
				return
			}
			// Log error if unmarshal fails, then proceed to fetch from DB
		}
	}

	// Cache miss or Redis is unavailable, fetch from database
	user := []models.User{}
	if result := config.DB.WithContext(ctx).Find(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users from database"})
		return
	}

	if config.RedisClient != nil {
		userJSON, err := json.Marshal(user)
		if err == nil {
			setCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			// Set cache in the background, log error if it fails but don't block the response
			go config.RedisClient.Set(setCtx, UserCacheKey, userJSON, UserCacheTTL)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"source": "database",
		"data":   user,
	})
}

func GetUserID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	userCacheKey := "user:" + id // สร้าง cache key เฉพาะสำหรับ user คนนี้

	// 1. ตรวจสอบใน Redis Cache ก่อน
	if config.RedisClient != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		cacheData, err := config.RedisClient.Get(cacheCtx, userCacheKey).Result()
		if err == nil { // ถ้าเจอใน cache (Cache hit)
			var user models.User
			if err := json.Unmarshal([]byte(cacheData), &user); err == nil {
				c.JSON(http.StatusOK, gin.H{"source": "cache", "data": user})
				return
			}
		}
	}

	var user models.User
	result := config.DB.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 3. นำข้อมูลที่ได้จากฐานข้อมูลไปเก็บใน Cache สำหรับการเรียกครั้งต่อไป
	if config.RedisClient != nil {
		userJSON, err := json.Marshal(user)
		if err == nil {
			setCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			// สั่งให้ set cache ทำงานเบื้องหลัง (background)
			go config.RedisClient.Set(setCtx, userCacheKey, userJSON, UserCacheTTL)
		}
	}

	// ส่งข้อมูลกลับไป
	c.JSON(http.StatusOK, gin.H{
		"source": "database",
		"data":   user,
	})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	result := config.DB.Where("?", c.Param("id")).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&user)
	if config.RedisClient != nil {
		// Invalidate cache
		go config.RedisClient.Del(c.Request.Context(), UserCacheKey)
		go config.RedisClient.Del(c.Request.Context(), "user:"+c.Param("id"))

		// Publish update event
		updateMsg, _ := json.Marshal(gin.H{"event": "user_updated", "user_id": user.Id})
		go config.RedisClient.Publish(c.Request.Context(), "user_updates", updateMsg)
	}

	c.JSON(http.StatusOK, &user)
}

func CreateUser(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{Username: input.Username, Name: input.Name, Email: input.Email, PasswordHash: string(hashedPassword)}
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Could not create user :" + result.Error.Error()})
		return
	}
	if config.RedisClient != nil {
		config.RedisClient.Del(c.Request.Context(), UserCacheKey)
	}
	c.JSON(http.StatusCreated, &user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	result := config.DB.Where("id = ?", c.Param("id")).Delete(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found or already deleted."})
		return
	}
	if config.RedisClient != nil {
		config.RedisClient.Del(c.Request.Context(), UserCacheKey)
	}
	c.JSON(http.StatusNoContent, nil)
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	var user models.User // สมมติว่านี่คือ User model ของคุณ

	// 1. Bind JSON body เข้ากับ struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. ค้นหา user จาก email ในฐานข้อมูล
	//    (โค้ดส่วนนี้ต้องปรับให้เข้ากับ DB logic ของคุณ)
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 3. เปรียบเทียบ password ที่ส่งมากับ hash ใน DB
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		// ถ้า password ไม่ตรงกัน
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 4. สร้าง JWT Token
	//    สร้าง claims สำหรับ token
	claims := jwt.MapClaims{
		"sub": user.Id,                               // Subject (user's ID)
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expiration time (24 hours)
		"iat": time.Now().Unix(),                     // Issued at
	}

	// สร้าง token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// เซ็น token ด้วย secret key ของคุณ (ควรเก็บไว้ใน environment variable)
	// ผมใช้ "YOUR_SECRET_KEY" เป็นตัวอย่าง คุณควรเปลี่ยนเป็นค่าของคุณเอง
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "YOUR_SECRET_KEY" // fallback สำหรับ local dev
	}

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// ส่ง token กลับไป
	c.JSON(http.StatusOK, gin.H{"token": t})
}
