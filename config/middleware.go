package config

import (
	"ental-health-system/models"
	"ental-health-system/utils"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/time/rate"
)

var (
	jwtSecret = []byte("your-secret-key") // 实际应用中应从环境变量获取
)

// Claims JWT声明结构
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// IPRateLimiter IP限流器
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter 创建新的IP限流器
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// GetLimiter 获取指定IP的限流器
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.GetLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// JWTMiddleware JWT认证中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "未授权，缺少认证令牌"})
			c.Abort()
			return
		}

		// 提取token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(401, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 从数据库查询token
		var tokenRecord models.Token
		if err := DB.Where("token = ? AND is_revoked = ? AND expires_at > ?",
			tokenString, false, time.Now()).First(&tokenRecord).Error; err != nil {
			c.JSON(401, gin.H{"error": "无效或已过期的认证令牌"})
			c.Abort()
			return
		}

		// 更新最后使用时间
		now := time.Now()
		DB.Model(&tokenRecord).Update("last_used", &now)

		// 验证JWT token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 记录日志信息
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		// 构建日志信息
		logInfo := map[string]interface{}{
			"status_code":  c.Writer.Status(),
			"latency_time": latencyTime,
			"client_ip":    c.ClientIP(),
			"req_method":   c.Request.Method,
			"req_uri":      c.Request.RequestURI,
			"user_id":      c.GetString("userID"),   // 从上下文获取用户ID
			"user_role":    c.GetString("userRole"), // 从上下文获取用户角色
		}

		// TODO: 将日志信息写入日志文件或发送到日志服务
		// 这里暂时只打印到控制台
		c.JSON(http.StatusOK, logInfo)
	}
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, role string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)

	claims := Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken 解析JWT令牌
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RoleAuthMiddleware 角色认证中间件
func RoleAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
			c.Abort()
			return
		}

		roleStr := role.(string)
		allowed := false
		for _, r := range allowedRoles {
			if r == roleStr {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问该资源"})
			c.Abort()
			return
		}

		c.Next()
	}
}
