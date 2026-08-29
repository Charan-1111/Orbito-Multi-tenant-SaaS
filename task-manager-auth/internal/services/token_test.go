package services

import (
	"task-manager-auth/internal/config"
	"task-manager-auth/internal/models"
	"task-manager-auth/internal/token"
	"testing"
)

func TestValidateTokenAcceptsRawAndBearerAuthorizationHeader(t *testing.T) {
	ts, err := token.NewTokenService("super-secret-key", "auth-service")
	if err != nil {
		t.Fatalf("failed to create token service: %v", err)
	}

	svc := &Service{
		config: &config.Configuration{Jwt: models.Jwt{
			AccessExpiryInMinutes:  15,
			RefreshExpiryInMinutes: 60,
		}},
		tokenService: ts,
	}

	accessToken, err := ts.GenerateToken("user-1", token.TokenTypeAccess, 15)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	if _, err = svc.ValidateToken(accessToken); err != nil {
		t.Fatalf("expected raw token to pass validation: %v", err)
	}

	if _, err = svc.ValidateToken("Bearer " + accessToken); err != nil {
		t.Fatalf("expected bearer token to pass validation: %v", err)
	}

	if _, err = svc.ValidateToken("Bearer"); err == nil {
		t.Fatal("expected malformed bearer header to fail")
	}
}

func TestRotateTokensRejectsAccessTokenAndAcceptsRefreshToken(t *testing.T) {
	ts, err := token.NewTokenService("super-secret-key", "auth-service")
	if err != nil {
		t.Fatalf("failed to create token service: %v", err)
	}

	svc := &Service{
		config: &config.Configuration{Jwt: models.Jwt{
			AccessExpiryInMinutes:  15,
			RefreshExpiryInMinutes: 60,
		}},
		tokenService: ts,
	}

	accessToken, err := ts.GenerateToken("user-1", token.TokenTypeAccess, 15)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	if _, _, err = svc.RotateTokens("Bearer " + accessToken); err == nil {
		t.Fatal("expected access token rotation to fail")
	}

	refreshToken, err := ts.GenerateToken("user-1", token.TokenTypeRefresh, 60)
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	newAccessToken, newRefreshToken, err := svc.RotateTokens("Bearer " + refreshToken)
	if err != nil {
		t.Fatalf("expected refresh token rotation to succeed: %v", err)
	}
	if newAccessToken == "" || newRefreshToken == "" {
		t.Fatal("rotation returned empty tokens")
	}
}
