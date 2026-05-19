package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(
	rdb *redis.Client,
	limit int,
	window time.Duration,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		key := "rate_limit:" + c.ClientIP()

		count, err := rdb.Incr(
			c.Request.Context(),
			key,
		).Result()

		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			rdb.Expire(
				c.Request.Context(),
				key,
				window,
			)
		}

		if count > int64(limit) {
			c.AbortWithStatusJSON(429, gin.H{
				"error": "Too many requests",
			})
			return
		}

		c.Next()
	}
}
