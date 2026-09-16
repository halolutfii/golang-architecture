package util

import (
	"context"
	"golang-clean-architecture/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
)

type TokenUtil struct {
	SecretKey string
	Redis     *redis.Client
}

func NewTokenUtil(secretKey string, redisClient *redis.Client) *TokenUtil {
	return &TokenUtil{
		SecretKey: secretKey,
		Redis:     redisClient,
	}
}

func (t TokenUtil) CreateToken(ctx context.Context, auth *model.Auth) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     auth.ID,
		"expire": time.Now().Add(time.Hour * 24 * 30).UnixMilli(),
	})

	jwtToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", nil
	}

	_, err = t.Redis.SetEx(ctx, jwtToken, auth.ID, time.Hour*24*30).Result()
	if err != nil {
		return "", nil
	}

	return jwtToken, nil
}

func (t TokenUtil) ParseToken(ctx context.Context, JwtToken string) (*model.Auth, error) {
	token, err := jwt.Parse(JwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	})

	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)

	expire := claims["expire"].(float64)
	if int64(expire) < time.Now().UnixMilli() {
		return nil, fiber.ErrUnauthorized
	}
	result, err := t.Redis.Exists(ctx, JwtToken).Result()
	if err != nil {
		return nil, err
	}

	if result == 0 {
		return nil, fiber.ErrUnauthorized
	}

	id := claims["id"].(string)
	auth := &model.Auth{
		ID: id,
	}

	return auth, nil
}
