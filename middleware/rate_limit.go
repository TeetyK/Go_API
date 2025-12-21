package middleware

import (
	"API/config"
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

		// ใช้ Pipeline เพื่อให้การทำงานของ INCR และ EXPIRE เกิดขึ้นพร้อมกัน (Atomic)
		var count int64
		pipe := config.RedisClient.Pipeline()
		incr := pipe.Incr(c.Request.Context(), key)
		// ตั้งค่า Expire ทุกครั้งที่เรียก เพื่อความแน่นอนและป้องกัน Race Condition
		pipe.Expire(c.Request.Context(), key, rateLimitPeriod)
		_, err := pipe.Exec(c.Request.Context())

		if err != nil {
			// ถ้า Redis มีปัญหา ก็ปล่อยผ่านไปก่อน แต่ควร log error ไว้
			c.Next()
			return
		}

		count = incr.Val()

		if count > rateLimitCount {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
