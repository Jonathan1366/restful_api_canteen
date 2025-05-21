package utils

import (
	"context"
	"fmt"
	entity "ubm-canteen/models"
	"google.golang.org/api/idtoken"
)

func VerifyGoogleIDToken(ctx context.Context, idToken string) (*entity.GoogleUser, error) {
    payload, err := idtoken.Validate(ctx, idToken, "")
    if err != nil {
        return nil, fmt.Errorf("failed to validate ID token: %w", err)
    }

    user := &entity.GoogleUser{
        Email:         payload.Claims["email"].(string),
        EmailVerified: payload.Claims["email_verified"].(bool),
        Name:          payload.Claims["name"].(string),
        Picture:       payload.Claims["picture"].(string),
        Sub:           payload.Subject,
    }

    return user, nil
}
