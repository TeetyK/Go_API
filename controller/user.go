package controller

import (
	"API/config"
	"API/models"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	UserCacheKey = "all_users"
	UserCacheTTL = 5 * time.Minute
)

func Paging(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		} else if pageSize > 100 {
			pageSize = 100
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "3"))
	ctx := c.Request.Context()
	user := []models.User{}
	var total int64
	db := config.DB.Model(&models.User{})
	if result := config.DB.WithContext(ctx).Find(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users from database"})
		return
	}
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not count users"})
		return
	}
	if err := db.Scopes(Paging(page, limit)).Find(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": user,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"limit":     limit,
			"last_page": int(math.Ceil(float64(total) / float64(limit))),
		},
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
	id := c.Param("id")

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&user)
	if config.RedisClient != nil {

		ctx := c.Request.Context()
		go config.RedisClient.Del(ctx, "user:"+id)

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
		// When a new user is created, we don't need to do anything to the cache
		// until that specific user is requested for the first time.
		// No need to invalidate "all_users".
	}
	c.JSON(http.StatusCreated, &user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	result := config.DB.Where("id = ?", c.Param("id")).Delete(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found or already deleted."})
		id := c.Param("id")

		// It's better to find the user first to ensure it exists.
		if err := config.DB.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
			return
		}

		// Now delete the user
		if err := config.DB.Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user."})
			return
		}
		if config.RedisClient != nil {
			ctx := c.Request.Context()
			go config.RedisClient.Del(ctx, "user:"+id)
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	var user models.User

	// 1. Bind JSON body เข้ากับ struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	claims := jwt.MapClaims{
		"sub": user.Id,                               // Subject (user's ID)
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expiration time (24 hours)
		"iat": time.Now().Unix(),                     // Issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")
	// if jwtSecret == "" {
	// 	jwtSecret = "YOUR_SECRET_KEY" // fallback สำหรับ local dev
	// }

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// ส่ง token กลับไป
	c.JSON(http.StatusOK, gin.H{"token": t})
}
