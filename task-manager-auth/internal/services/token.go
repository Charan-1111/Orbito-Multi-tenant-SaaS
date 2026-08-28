package services

import "task-manager-auth/internal/token"

func (s *Service) ValidateToken(token string) (*token.Claims, error){
	tokenClaims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return tokenClaims, nil
}