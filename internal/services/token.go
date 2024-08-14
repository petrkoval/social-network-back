package services

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/petrkoval/social-network-back/internal/config"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/rs/zerolog"
	"time"
)

type TokenStorage interface {
	Find(ctx context.Context, refreshToken string) (*domain.Token, error)
	Save(ctx context.Context, token domain.Token) error
	Delete(ctx context.Context, refreshToken string) error
}

type TokenService struct {
	storage *TokenStorage
	logger  *zerolog.Logger
	cfg     *config.TokensConfig
}

func NewTokenService(s *TokenStorage, l *zerolog.Logger, cfg *config.TokensConfig) *TokenService {
	return &TokenService{
		storage: s,
		logger:  l,
		cfg:     cfg,
	}
}

func (s *TokenService) GenerateTokens(user domain.User) (accessToken, refreshToken string, err error) {
	var (
		accessTokenClaims    domain.TokenClaims
		refreshTokenClaims   domain.TokenClaims
		unsignedAccessToken  *jwt.Token
		unsignedRefreshToken *jwt.Token
	)
	s.logger.Debug().Msg("generating tokens")

	accessTokenClaims = domain.TokenClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "token service",
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		},
	}

	refreshTokenClaims = domain.TokenClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "token service",
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}

	unsignedAccessToken = jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	unsignedRefreshToken = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	accessToken, err = unsignedAccessToken.SignedString([]byte(s.cfg.AccessSecret))
	if err != nil {
		return "", "", jwtSigningErr
	}

	refreshToken, err = unsignedRefreshToken.SignedString([]byte(s.cfg.RefreshSecret))
	if err != nil {
		return "", "", jwtSigningErr
	}

	s.logger.Debug().Msg("tokens generated")
	return accessToken, refreshToken, nil
}

func (s *TokenService) VerifyAccessToken(accessToken string) (*domain.User, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.ParseWithClaims(accessToken, &domain.TokenClaims{}, s.validateAccessSigningMethod)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.TokenClaims); ok && token.Valid {
		if claims.ExpiresAt.Before(time.Now()) {
			return nil, tokenExpiredErr
		}

		return &domain.User{
			ID:       claims.Subject,
			Username: claims.Username,
		}, nil
	}

	return nil, invalidTokenErr
}

func (s *TokenService) VerifyRefreshToken(refreshToken string) (*domain.User, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.ParseWithClaims(refreshToken, &domain.TokenClaims{}, s.validateRefreshSigningMethod)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.TokenClaims); ok && token.Valid {
		if claims.ExpiresAt.Before(time.Now()) {
			return nil, tokenExpiredErr
		}

		return &domain.User{
			ID:       claims.Subject,
			Username: claims.Username,
		}, nil
	}

	return nil, invalidTokenErr
}

func (s *TokenService) validateAccessSigningMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, unexpectedSigningMethodErr
	}

	return []byte(s.cfg.AccessSecret), nil
}

func (s *TokenService) validateRefreshSigningMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, unexpectedSigningMethodErr
	}

	return []byte(s.cfg.RefreshSecret), nil
}
