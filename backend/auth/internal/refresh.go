package internal

import (
	"auth/configs"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const refreshTokenPrefix = "refresh_token:"

func GenerateRefreshToken(email string) (string, error) {
	refreshToken := uuid.New().String()
	key := refreshTokenPrefix + refreshToken
	
	ctx := context.Background()
	ttl := time.Duration(configs.RefreshTokenTTLDays) * 24 * time.Hour
	
	err := configs.RedisClient.Set(ctx, key, email, ttl).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}
	
	return refreshToken, nil
}

func ValidateRefreshToken(refreshToken string) (string, error) {
	key := refreshTokenPrefix + refreshToken
	ctx := context.Background()
	
	email, err := configs.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("invalid or expired refresh token")
	}
	
	return email, nil
}

func RevokeRefreshToken(refreshToken string) error {
	key := refreshTokenPrefix + refreshToken
	ctx := context.Background()
	
	err := configs.RedisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	
	return nil
}