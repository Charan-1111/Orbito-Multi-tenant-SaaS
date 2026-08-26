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

	fmt.Println(username,  " ", hashedPassword)

	// now should we store the username and password in the database..

	

	return nil
}
