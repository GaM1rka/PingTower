package configs

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	Logger               *log.Logger
	JWTSecret            string
	ServerPort           string
	JWTExpiry            string
	RedisClient          *redis.Client
	RedisHost            string
	RedisPort            string
	RedisPassword        string
	RefreshTokenTTLDays  int
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func InitConfig() {
	JWTSecret = getEnv("JWT_SECRET", "hackathon_secret_key_2024")
	ServerPort = getEnv("SERVER_PORT", "8081")
	JWTExpiry = getEnv("JWT_EXPIRY_HOURS", "1")
	RedisHost = getEnv("REDIS_HOST", "localhost")
	RedisPort = getEnv("REDIS_PORT", "6379")
	RedisPassword = getEnv("REDIS_PASSWORD", "")
	
	ttlStr := getEnv("REFRESH_TOKEN_TTL_DAYS", "7")
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		ttl = 7
	}
	RefreshTokenTTLDays = ttl
}

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     RedisHost + ":" + RedisPort,
		Password: RedisPassword,
		DB:       0,
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	Logger.Println("Redis connected successfully")
	return nil
}

func InitLogger() {
	Logger = log.New(os.Stdout, "AUTH: ", log.LstdFlags)
}