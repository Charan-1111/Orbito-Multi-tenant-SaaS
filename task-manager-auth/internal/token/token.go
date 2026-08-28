package token

import (
	"task-manager-auth/internal/constants"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`

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

func (t *TokenService) GenerateToken(userId string, expiry int) (string, error) {
	now := time.Now()
	claims := Claims{
		UserId: userId,
		Email:  userId,
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
	accessToken, err := t.GenerateToken(userId, accessExpiry)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := t.GenerateToken(userId, refreshExpiry)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// TODO : Token validation and parsing logic needs to be implemented below
