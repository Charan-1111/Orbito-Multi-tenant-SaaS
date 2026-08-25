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

	user := strings.Split(strings.TrimSpace(userDetails), " ")
	username := user[0]
	passWord := user[1]

	fmt.Println("Request ID:", requestId)
	fmt.Println("Username:", username)
	fmt.Println("Password:", passWord)
	return nil
}
