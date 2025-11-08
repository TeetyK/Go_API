package controller

import (
	"API/config"
	"API/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	user := []models.User{}
	config.DB.Find(&user)
	c.JSON(http.StatusOK, user)
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
	c.JSON(http.StatusCreated, &user)
}
func DeleteUser(c *gin.Context) {
	var user models.User
	result := config.DB.Where("id = ?", c.Param("id")).Delete(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found or already deleted."})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
