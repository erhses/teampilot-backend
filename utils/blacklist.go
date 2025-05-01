package utils

// import "sync"

// var (
// 	blacklistedTokens = make(map[string]bool)
// 	mu                sync.RWMutex
// )

// func BlacklistToken(token string) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	blacklistedTokens[token] = true
// }

// func IsTokenBlacklisted(token string) bool {
// 	mu.RLock()
// 	defer mu.RUnlock()
// 	return blacklistedTokens[token]
// }


import (
	"context"
	"teampilot/integrations/rdb"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	redis *redis.Client
	mu 	sync.RWMutex
}

var (
	tokenManager     *TokenManager
	tokenManagerOnce sync.Once
)

func GetTokenManager() *TokenManager {
    tokenManagerOnce.Do(func() {
        redisClient := rdb.GetRedisClient()
        tokenManager = &TokenManager{
            redis: redisClient,
        }
    })
    return tokenManager
}


func (tm *TokenManager) BlacklistToken(ctx context.Context, tokenString string) error {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return errors.New("invalid JWT format")
	}
	signature := parts[2]

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return GetAccessSecretKey(), nil
	})

	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid token expiration")
	}

	expTime := time.Unix(int64(exp), 0)
	ttl := time.Until(expTime)

	if ttl <= 0 {
		return nil
	}

	key := "blacklist:" + signature

	tm.mu.Lock()
	defer tm.mu.Unlock()

	err = tm.redis.Set(ctx, key, "revoked", ttl).Err()
	if err != nil {
		return err
	}

	return nil
}


func (tm *TokenManager) IsTokenBlacklisted(ctx context.Context, tokenString string) bool {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return true
	}
	signature := parts[2]

	key := "blacklist:" + signature
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	exists, err := tm.redis.Exists(ctx, key).Result()
	if err != nil {
		return true
	}
	return exists > 0
}


func (tm *TokenManager) ValidateToken(ctx context.Context, tokenString string) (jwt.MapClaims, error) {
	if tm.IsTokenBlacklisted(ctx, tokenString) {
		return nil, errors.New("token has been revoked")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return GetAccessSecretKey(), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("Таны нэвтрэх хугацаа дууссан байна")
	}

	return claims, nil
}

// func (tm *TokenManager) CleanupBlacklist(ctx context.Context) error {
// 	iter := tm.redis.Scan(ctx, 0, "blacklist:*", 100).Iterator()
// 	for iter.Next(ctx) {
// 		key := iter.Val()
// 		ttl := tm.redis.TTL(ctx, key).Val()
// 		if ttl <= 0 {
// 			tm.redis.Del(ctx, key)
// 		}
// 	}
// 	return iter.Err()
// }

