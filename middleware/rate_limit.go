// ในไฟล์ middleware/rate_limiter.go (ไฟล์ใหม่)
package middleware

import (
	"API/config" // Import config ของคุณ
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	rateLimitPeriod = 1 * time.Minute
	rateLimitCount  = 5 // อนุญาต 5 ครั้งต่อนาที
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.RedisClient == nil {
			c.Next()
			return
		}

		// ใช้ IP Address เป็น key
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		// เพิ่มค่าใน key นี้, ถ้า key ไม่มีอยู่จะถูกสร้างและมีค่าเป็น 1
		count, err := config.RedisClient.Incr(c.Request.Context(), key).Result()
		if err != nil {
			// ถ้า Redis มีปัญหา ก็ปล่อยผ่านไปก่อน แต่ควร log error ไว้
			c.Next()
			return
		}

		// ถ้าเป็นการสร้าง key ครั้งแรก (count == 1) ให้ตั้งเวลาหมดอายุ
		if count == 1 {
			config.RedisClient.Expire(c.Request.Context(), key, rateLimitPeriod)
		}

		if count > rateLimitCount {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
