package controller

import (
	"API/config"
	"API/models"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	AllProductsCacheKey = "all_products"
	ProductCacheTTL     = 5 * time.Minute
)

// GetProducts godoc
// @Summary Get all products
// @Description Get a list of all products, with caching.
// @Tags products
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
func GetProducts(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Try to get from cache first
	if config.RedisClient != nil {
		cacheData, err := config.RedisClient.Get(ctx, AllProductsCacheKey).Result()
		if err == nil {
			var products []models.Product
			if json.Unmarshal([]byte(cacheData), &products) == nil {
				c.JSON(http.StatusOK, gin.H{"source": "cache", "data": products})
				return
			}
		}
	}

	// 2. If cache miss, get from DB
	var products []models.Product
	if result := config.DB.WithContext(ctx).Find(&products); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
		return
	}

	// 3. Set to cache for next time (in background)
	if config.RedisClient != nil {
		productsJSON, err := json.Marshal(products)
		if err == nil {
			go config.RedisClient.Set(context.Background(), AllProductsCacheKey, productsJSON, ProductCacheTTL)
		}
	}

	c.JSON(http.StatusOK, gin.H{"source": "database", "data": products})
}

// GetProductByID godoc
// @Summary Get a single product by its ID
// @Description Get detailed information for a specific product.
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Router /products/{id} [get]
func GetProductByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	productCacheKey := "product:" + id

	// 1. Try to get from cache
	if config.RedisClient != nil {
		cachedProduct, err := config.RedisClient.Get(ctx, productCacheKey).Result()
		if err == nil {
			var product models.Product
			if json.Unmarshal([]byte(cachedProduct), &product) == nil {
				c.JSON(http.StatusOK, gin.H{"source": "cache", "data": product})
				return
			}
		}
	}

	// 2. If cache miss, get from DB
	var product models.Product
	if result := config.DB.WithContext(ctx).First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// 3. Set to cache
	if config.RedisClient != nil {
		productJSON, err := json.Marshal(product)
		if err == nil {
			go config.RedisClient.Set(context.Background(), productCacheKey, productJSON, ProductCacheTTL)
		}
	}

	c.JSON(http.StatusOK, gin.H{"source": "database", "data": product})
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Adds a new product to the database.
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product object"
// @Success 201 {object} models.Product
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := config.DB.Create(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create product: " + result.Error.Error()})
		return
	}

	// Invalidate cache for all products list
	if config.RedisClient != nil {
		go config.RedisClient.Del(context.Background(), AllProductsCacheKey)
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Updates a product's details by its ID.
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product object"
// @Success 200 {object} models.Product
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if result := config.DB.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&product)

	// Invalidate caches
	if config.RedisClient != nil {
		productCacheKey := "product:" + id
		go config.RedisClient.Del(context.Background(), AllProductsCacheKey)
		go config.RedisClient.Del(context.Background(), productCacheKey)
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Deletes a product by its ID.
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 204 "No Content"
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Convert id to uint for GORM
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	result := config.DB.Delete(&models.Product{}, uid)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Invalidate caches
	if config.RedisClient != nil {
		productCacheKey := "product:" + id
		go config.RedisClient.Del(context.Background(), AllProductsCacheKey)
		go config.RedisClient.Del(context.Background(), productCacheKey)
	}

	c.Status(http.StatusNoContent)
}