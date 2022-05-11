package middleware

import (
	"errors"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func JWTAuth(c *gin.Context) {
	token := c.Request.Header.Get("x-token")
	zap.S().Info("x-token:", token)
	if token == "" {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"msg": "请登录",
		})
		c.Abort()
		return
	}
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		if err == TokenExpired {
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "授权已过期",
				})
				c.Abort()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "未登录",
			"msg":    err.Error(),
		})
		c.Abort()
		return
	}
	c.Set("claims", claims)
	c.Set("userId", claims.ID)
	c.Next()
}

type JWT struct {
	SigningKey []byte // 签名密钥
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

// NewJWT 返回JWT指针
func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServiceConfig.JWTInfo.SigningKey), //可以设置过期时间
	}
}

func (j *JWT) CreateToken(claims model.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (*model.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
