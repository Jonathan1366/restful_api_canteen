package utils

import (
	"context"
	"fmt"
	entity "ubm-canteen/models"

	"google.golang.org/api/idtoken"
)

func VerifyGoogleIDToken(ctx context.Context, idToken string, audience string) (*entity.GoogleUser, error) {
    
    payload, err := idtoken.Validate(ctx, idToken, audience)
    
    if err != nil {
        return nil, fmt.Errorf("failed to validate ID token: %w", err)
    }
    
    email, ok := payload.Claims["email"].(string)
    if !ok{
        return nil, fmt.Errorf("email claim not found in ID token")
    }
    
    emailVerified, ok := payload.Claims["email_verified"].(bool)
    
    if !ok{
        return nil, fmt.Errorf("email_verified claim not found in ID token")
    }
    
    name, _ := payload.Claims["name"].(string)
    picture, _ := payload.Claims["picture"].(string)

    return &entity.GoogleUser{
        Email:         email,
        EmailVerified: emailVerified,
        Name:          name,
        Picture:      picture,
        Sub: payload.Subject,      
    }, nil
}
