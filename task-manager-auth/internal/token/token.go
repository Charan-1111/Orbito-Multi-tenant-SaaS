package token

import (
	"fmt"
	"task-manager-auth/internal/constants"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Claims struct {
	UserId    string `json:"userId"`
	Email     string `json:"email"`
	TokenType string `json:"tokenType"`

	jwt.RegisteredClaims
}

type TokenService struct {
	secret []byte
	issuer string
}

func NewTokenService(secret, issuer string) (*TokenService, error) {
	if secret == "" {
		return nil, constants.ErrMissingJwtSecret
	}

	return &TokenService{
		secret: []byte(secret),
		issuer: issuer,
	}, nil
}

func (t *TokenService) GenerateToken(userId, tokenType string, expiry int) (string, error) {
	now := time.Now()
	claims := Claims{
		UserId:    userId,
		Email:     userId,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.issuer,
			Subject:   userId,
			Audience:  jwt.ClaimStrings{"my-api"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiry) * time.Minute)),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	// generating the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signing the token using secret key
	signedToken, err := token.SignedString(t.secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (t *TokenService) GenerateTokens(userId string, accessExpiry, refreshExpiry int) (string, string, error) {
	accessToken, err := t.GenerateToken(userId, TokenTypeAccess, accessExpiry)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := t.GenerateToken(userId, TokenTypeRefresh, refreshExpiry)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (t *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf(
					"%w : unexpected signing method %s",
					constants.ErrInvalidToken,
					token.Method.Alg(),
				)
			}

			return t.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(t.issuer),
		jwt.WithAudience("my-api"),
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", constants.ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, constants.ErrInvalidToken
	}

	return claims, nil
}
