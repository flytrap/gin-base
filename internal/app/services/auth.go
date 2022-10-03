package services

import (
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/flytrap/gin-base/internal/app/config"
	"github.com/flytrap/gin-base/pkg/redis"
)

const defaultKey = "flytrap"

var defaultOptions = AuthOption{
	TokenType:     "Bearer",
	Expired:       7200,
	SigningMethod: jwt.SigningMethodHS512,
	SigningKey:    []byte(defaultKey),
	KeyFunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrHashUnavailable
		}
		return []byte(defaultKey), nil
	},
}

type AuthService interface {
	ParseUserID(tokenString string) (string, error)
	GetToken(string) string
}

type AuthOption struct {
	SigningMethod jwt.SigningMethod
	SigningKey    interface{}
	KeyFunc       jwt.Keyfunc
	Expired       int
	TokenType     string
}

func NewAuthService() AuthService {
	cfg := config.C.JWTAuth
	var method jwt.SigningMethod
	switch cfg.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}
	opt := AuthOption{
		Expired:    cfg.Expired,
		SigningKey: cfg.SigningKey,
		KeyFunc: func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKey
			}
			return []byte(cfg.SigningKey), nil
		},
		SigningMethod: method,
		TokenType:     "Bearer",
	}

	a := InitAuth(&opt)
	return a
}

// New 创建认证实例
func InitAuth(opt *AuthOption) AuthService {
	o := defaultOptions
	if opt != nil {
		return &JwtAuthService{Opts: opt}
	}
	return &JwtAuthService{Opts: &o}
}

type JwtAuthService struct {
	Opts  *AuthOption
	Store *redis.Store
}

// 解析令牌
func (a *JwtAuthService) parseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, a.Opts.KeyFunc)
	if err != nil || !token.Valid {
		return nil, jwt.ErrHashUnavailable
	}

	return token.Claims.(*jwt.StandardClaims), nil
}

// ParseUserID 解析用户ID
func (a *JwtAuthService) ParseUserID(tokenString string) (string, error) {
	if tokenString == "" {
		return "", jwt.ErrHashUnavailable
	}

	claims, err := a.parseToken(tokenString)
	if err != nil {
		return "", err
	}
	err = claims.Valid()
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

func (a *JwtAuthService) GetToken(auth string) string {
	var token string
	prefix := a.Opts.TokenType
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = strings.TrimSpace(auth[len(prefix):])
	}
	return token
}
