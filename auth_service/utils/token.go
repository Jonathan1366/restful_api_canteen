package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

func VerifyGoogleIDToken(ctx context.Context, idToken string) (*idtoken.Payload, error) {
    payload, err := idtoken.Validate(ctx, idToken, os.Getenv("WEB_CLIENT_ID"))
    if err != nil {
        return nil, fmt.Errorf("failed to validate ID token: %w", err)
    }
    return payload, nil
}

func ExchangeAuthCode(ctx context.Context, code string) (*oauth2.Token, error) {
    cfg := &oauth2.Config{
        ClientID:    os.Getenv("WEB_CLIENT_ID"),
        ClientSecret: os.Getenv("WEB_CLIENT_SECRET"),
        Endpoint: google.Endpoint,
    }
    token, err:= cfg.Exchange(ctx, code)
    if err != nil {
        return nil, fmt.Errorf("exchange failed: %w", err)
    }
    return token, nil
}

func StoreRefreshToken(ctx context.Context, rdb *redis.Client, userID, refreshToken string) error {
  h := sha256.Sum256([]byte(refreshToken))
  hash := hex.EncodeToString(h[:])
  key := fmt.Sprintf("refresh:user:%s", userID)
  return rdb.Set(ctx, key, hash, 30*24*time.Hour).Err()
}
