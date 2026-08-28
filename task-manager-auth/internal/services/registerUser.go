package services

import (
	"context"
	"fmt"
	"strings"
	"task-manager-auth/internal/utils"
)

func (s *Service) RegisterUser(ctx context.Context, requestId, details string) error {
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

func (s *Service) UserLogin(ctx context.Context, requestId, details string) error {
	userDetails, err := utils.Base64Decode(details)
	if err != nil {
		return fmt.Errorf("Decoding the user detalils failed : %w", err)
	}

	user := strings.Split(strings.TrimSpace(userDetails), ":")
	username := user[0]
	passWord := user[1]

	hashedPassword, err := utils.HashPassword(passWord)
	if err != nil {
		return fmt.Errorf("Hashing the password : %w", err)
	}

	userCnt, err := s.database.CheckUserExistance(ctx, username, hashedPassword)
	if err != nil || userCnt == 0 {
		return fmt.Errorf("Invalid Credentials")
	}

	return nil
}
