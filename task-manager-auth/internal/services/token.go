package services

import (
	"strings"
	"task-manager-auth/internal/constants"
	"task-manager-auth/internal/token"
)

func (s *Service) ValidateToken(authHeader string) (*token.Claims, error) {
	parts := strings.Fields(strings.TrimSpace(authHeader))
	if len(parts) == 0 {
		return nil, constants.ErrInvalidToken
	}

	tokenValue := parts[0]
	if len(parts) == 2 {
		if !strings.EqualFold(parts[0], "Bearer") {
			return nil, constants.ErrInvalidToken
		}
		tokenValue = parts[1]
	}
	if len(parts) > 2 || tokenValue == "" {
		return nil, constants.ErrInvalidToken
	}

	tokenClaims, err := s.tokenService.ValidateToken(tokenValue)
	if err != nil {
		return nil, err
	}

	return tokenClaims, nil
}

func (s *Service) RotateTokens(refreshToken string) (string, string, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	if claims.TokenType != token.TokenTypeRefresh {
		return "", "", constants.ErrInvalidRefreshToken
	}

	accessToken, rotatedRefreshToken, err := s.tokenService.GenerateTokens(claims.UserId, s.config.Jwt.AccessExpiryInMinutes, s.config.Jwt.RefreshExpiryInMinutes)
	if err != nil {
		return "", "", err
	}

	return accessToken, rotatedRefreshToken, nil
}
