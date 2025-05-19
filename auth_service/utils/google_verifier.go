package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	entity "ubm-canteen/models"
)

func VerifyGoogleIDToken(ctx context.Context, idToken string) (*entity.GoogleUser, error){
	resp, err:= http.Get("https://oauth2.googleapis.com/tokeninfo?id_token="+idToken )
	if err != nil {
		return nil, fmt.Errorf("failed to verify idToken: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		return nil, errors.New("invalid Google ID Token")
	}

	var user entity.GoogleUser
	if err:= json.NewDecoder(resp.Body).Decode(&user); err!=nil{
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if user.Email=="" || user.Sub ==""{
		return nil, errors.New("missing email or sub from token")
	}
	return &user, nil
} 