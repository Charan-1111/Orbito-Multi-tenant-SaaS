package services

import (
	"context"
	"fmt"
	"strings"
	"task-manager-auth/internal/utils"
)

func (s *Service) RegisterUser(ctx context.Context, requestId, details string) error {
	// TODO : Need to include the logic to generate the userId for the given user

	userDetails, err := utils.Base64Decode(details)
	if err != nil {
		return fmt.Errorf("Decoding the user details failed: %w", err)
	}

	user := strings.Split(strings.TrimSpace(userDetails), ":")
	username := user[0]
	passWord := user[1]

	hashedPassword, err := utils.HashPassword(passWord)
	if err != nil {
		return fmt.Errorf("Hashing the password : %w", err)
	}

	err = s.database.RegisterUser(ctx, username, hashedPassword)
	if err != nil {
		return fmt.Errorf("Registering User : %w", err)
	}

	return nil
}

func (s *Service) UserLogin(ctx context.Context, requestId, details string) (string, string, error) {
	userDetails, err := utils.Base64Decode(details)
	if err != nil {
		return "", "", fmt.Errorf("Decoding the user detalils failed : %w", err)
	}

	user := strings.Split(strings.TrimSpace(userDetails), ":")
	username := user[0]
	passWord := user[1]

	passwordHash, err := s.database.CheckUserExistance(ctx, username)
	if err != nil || utils.VerifyPassword(passWord, passwordHash) != nil {
		return "", "", fmt.Errorf("Invalid Credentials")
	}

	// After successful login we are going to generate the required tokens
	accessToken, refreshToken, err := s.tokenService.GenerateTokens(username, s.config.Jwt.AccessExpiryInMinutes, s.config.Jwt.RefreshExpiryInMinutes)

	return accessToken, refreshToken, nil
}
