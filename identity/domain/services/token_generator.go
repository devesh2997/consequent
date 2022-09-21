package services

import (
	"context"
	"fmt"
	"time"

	"github.com/devesh2997/consequent/errorx"
	"github.com/devesh2997/consequent/identity/constants"
	"github.com/devesh2997/consequent/identity/domain/entities"
	"github.com/devesh2997/consequent/identity/domain/repositories"
	userEntities "github.com/devesh2997/consequent/user/domain/entities"
	"github.com/golang-jwt/jwt"
)

const (
	jwtExpiryDuration          = time.Minute * 10
	refreshTokenExpiryDuration = time.Hour * 24
)

type TokenService interface {
	Generate(ctx context.Context, user userEntities.User) (*entities.Token, error)
}

func NewTokenService(repo repositories.TokenRepo) TokenService {
	return tokenService{repo: repo}
}

type tokenService struct {
	repo repositories.TokenRepo
}

func (service tokenService) Generate(ctx context.Context, user userEntities.User) (*entities.Token, error) {
	now := time.Now().UTC()
	jwtExpiryAt := now.Add(jwtExpiryDuration)
	refreshTokenExpiryAt := now.Add(refreshTokenExpiryDuration)

	jwtClaims := service.getJWTClaims(user, jwtExpiryAt.Unix())
	jwtTokenStr, err := service.signClaims(jwtClaims)
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := service.getRefreshTokenClaims(user.ID, refreshTokenExpiryAt.Unix())
	refreshTokenStr, err := service.signClaims(refreshTokenClaims)
	if err != nil {
		return nil, err
	}

	refreshToken := entities.RefreshToken{
		Token:     refreshTokenStr,
		Status:    constants.REFRESH_TOKEN_STATUS_ACTIVE,
		CreatedAt: time.Now(),
		ExpiryAt:  refreshTokenExpiryAt,
		UpdatedAt: time.Now(),
	}

	if err := service.repo.SaveRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	token := entities.Token{
		JWT: entities.JWT{
			Token:    jwtTokenStr,
			ExpiryAt: jwtExpiryAt,
		},
		RefreshToken: refreshToken,
	}

	return &token, nil
}

func (service tokenService) getJWTClaims(user userEntities.User, exp int64) jwt.MapClaims {
	claims := make(jwt.MapClaims)
	claims["sub"] = user.ID                     // Subject of the token (i.e. the user)
	claims["usr"] = service.getJWTPayload(user) // User data.
	claims["exp"] = exp                         // The expiration time after which the token must be disregarded.
	claims["iat"] = time.Now()                  // The time at which the token was issued.
	claims["nbf"] = time.Now()                  // The time before which the token must be disregarded.

	return claims
}

func (service tokenService) getRefreshTokenClaims(userID int64, exp int64) jwt.MapClaims {
	claims := make(jwt.MapClaims)
	claims["sub"] = userID     // Subject of the token (i.e. the user)
	claims["exp"] = exp        // The expiration time after which the token must be disregarded.
	claims["iat"] = time.Now() // The time at which the token was issued.
	claims["nbf"] = time.Now() // The time before which the token must be disregarded.

	return claims
}

func (service tokenService) signClaims(claims jwt.MapClaims) (string, error) {
	privateKey, err := service.repo.GetPrivateKey()
	if err != nil {
		return "", errorx.NewSystemError(-1, err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", errorx.NewSystemError(-1, fmt.Errorf("create: parse key: %w", err))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "1" // TODO (devesh2997) | this should not be hardcoded. a separate key manager service is advised.
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "nil", errorx.NewSystemError(-1, fmt.Errorf("create: sign token: %w", err))
	}

	return tokenStr, nil
}

func (service tokenService) getJWTPayload(user userEntities.User) map[string]interface{} {
	return map[string]interface{}{
		"id":     user.ID,
		"email":  user.Email,
		"mobile": user.Mobile,
		"name":   user.Name,
		"gender": user.Gender,
	}
}
