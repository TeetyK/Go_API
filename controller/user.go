package controller

import (
	"API/config"
	"API/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	UserCacheKey = "all_users"
	UserCacheTTL = 5 * time.Minute
)

func GetUsers(c *gin.Context) {
	if config.RedisClient != nil {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		cacheData, err := config.RedisClient.Get(ctx, UserCacheKey).Result()
		if err == nil {
			user := []models.User{}
			if err := json.Unmarshal([]byte(cacheData), &user); err == nil {
				c.JSON(http.StatusOK, gin.H{"source": "cache", "data": user})
			}
			return
		}
	}
	user := []models.User{}
	config.DB.Find(&user)
	if config.RedisClient != nil {
		userJSON, err := json.Marshal(user)
		if err == nil {
			config.RedisClient.Set(context.Background(), UserCacheKey, userJSON, UserCacheTTL)
		}
	}
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
	config.RedisClient.Del(context.Background(), UserCacheKey)
	c.JSON(http.StatusOK, &user)
}
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Could not create user :" + result.Error.Error()})
		return
	}
	config.RedisClient.Del(context.Background(), UserCacheKey)
	c.JSON(http.StatusCreated, &user)
}
func DeleteUser(c *gin.Context) {
	var user models.User
	result := config.DB.Where("id = ?", c.Param("id")).Delete(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found or already deleted."})
		return
	}
	config.RedisClient.Del(context.Background(), UserCacheKey)
	c.JSON(http.StatusNoContent, nil)
}
