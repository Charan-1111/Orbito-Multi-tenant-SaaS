package services

import (
	"context"
	"fmt"
	"strings"
	"task-manager-auth/internal/utils"
)

func (s *Service) RegisterUser(ctx context.Context, requestId, details string) error {
	logger := s.log.WithRequestID(requestId)
	logger.Info().Str("event", "register_user").Msg("processing registration request")

	userDetails, err := utils.Base64Decode(details)
	if err != nil {
		logger.Error().Err(err).Str("event", "register_user").Msg("failed to decode registration payload")
		return fmt.Errorf("Decoding the user details failed: %w", err)
	}

	user := strings.Split(strings.TrimSpace(userDetails), ":")
	username := user[0]
	passWord := user[1]

	hashedPassword, err := utils.HashPassword(passWord)
	if err != nil {
		logger.Error().Err(err).Str("event", "register_user").Msg("failed to hash password")
		return fmt.Errorf("Hashing the password : %w", err)
	}

	err = s.database.RegisterUser(ctx, username, hashedPassword)
	if err != nil {
		logger.Error().Err(err).Str("event", "register_user").Msg("failed to persist user")
		return fmt.Errorf("Registering User : %w", err)
	}

	logger.Info().Str("event", "register_user").Str("user", username).Msg("registration request completed")
	return nil
}

func (s *Service) UserLogin(ctx context.Context, requestId, details string) (string, string, error) {
	logger := s.log.WithRequestID(requestId)
	logger.Info().Str("event", "login").Msg("processing login request")

	userDetails, err := utils.Base64Decode(details)
	if err != nil {
		logger.Error().Err(err).Str("event", "login").Msg("failed to decode login payload")
		return "", "", fmt.Errorf("Decoding the user detalils failed : %w", err)
	}

	user := strings.Split(strings.TrimSpace(userDetails), ":")
	username := user[0]
	passWord := user[1]

	passwordHash, err := s.database.CheckUserExistance(ctx, username)
	if err != nil || utils.VerifyPassword(passWord, passwordHash) != nil {
		logger.Error().Err(err).Str("event", "login").Str("user", username).Msg("invalid credentials for login request")
		return "", "", fmt.Errorf("Invalid Credentials")
	}

	accessToken, refreshToken, err := s.tokenService.GenerateTokens(username, s.config.Jwt.AccessExpiryInMinutes, s.config.Jwt.RefreshExpiryInMinutes)
	if err != nil {
		logger.Error().Err(err).Str("event", "login").Str("user", username).Msg("failed to generate tokens")
		return "", "", err
	}

	logger.Info().Str("event", "login").Str("user", username).Msg("login request completed")
	return accessToken, refreshToken, nil
}
