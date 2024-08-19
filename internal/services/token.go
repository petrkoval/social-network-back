package services

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/petrkoval/social-network-back/internal/config"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"time"
)

type TokenStorage interface {
	Find(ctx context.Context, refreshToken string) (*domain.Token, error)
	Save(ctx context.Context, token domain.Token) error
	Delete(ctx context.Context, refreshToken string) error
}

type TokenService struct {
	Storage TokenStorage
	logger  *zerolog.Logger
	cfg     *config.TokensConfig
}

func NewTokenService(s TokenStorage, l *zerolog.Logger, cfg *config.TokensConfig) *TokenService {
	return &TokenService{
		Storage: s,
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
		return "", "", errors.Wrap(JwtSigningErr, "TokenService.GenerateTokens")
	}

	refreshToken, err = unsignedRefreshToken.SignedString([]byte(s.cfg.RefreshSecret))
	if err != nil {
		return "", "", errors.Wrap(JwtSigningErr, "TokenService.GenerateTokens")
	}

	s.logger.Debug().
		Str("accessToken", accessToken).
		Str("refreshToken", refreshToken).
		Msg("tokens generated")

	return accessToken, refreshToken, nil
}

func (s *TokenService) VerifyAccessToken(accessToken string) (*domain.User, error) {
	var (
		token *jwt.Token
		err   error
	)
	s.logger.Debug().Msg("verifying access token")

	token, err = jwt.ParseWithClaims(accessToken, &domain.TokenClaims{}, s.validateAccessSigningMethod)
	if err != nil {
		return nil, errors.Wrap(err, "TokenService.VerifyAccessToken")
	}

	if claims, ok := token.Claims.(*domain.TokenClaims); ok && token.Valid {
		if claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.Wrap(TokenExpiredErr, "TokenService.VerifyAccessToken")
		}

		entity := &domain.User{
			ID:       claims.Subject,
			Username: claims.Username,
		}

		s.logger.Debug().
			Str("userID", entity.ID).
			Str("Username", entity.Username).
			Msg("access token verified")

		return entity, nil
	}

	return nil, errors.Wrap(InvalidTokenErr, "TokenService.VerifyAccessToken")
}

func (s *TokenService) VerifyRefreshToken(refreshToken string) (*domain.User, error) {
	var (
		token *jwt.Token
		err   error
	)
	s.logger.Debug().Msg("verifying refresh token")

	token, err = jwt.ParseWithClaims(refreshToken, &domain.TokenClaims{}, s.validateRefreshSigningMethod)
	if err != nil {
		return nil, errors.Wrap(err, "TokenService.VerifyRefreshToken")
	}

	if claims, ok := token.Claims.(*domain.TokenClaims); ok && token.Valid {
		if claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.Wrap(TokenExpiredErr, "TokenService.VerifyRefreshToken")
		}

		entity := &domain.User{
			ID:       claims.Subject,
			Username: claims.Username,
		}

		s.logger.Debug().
			Str("userID", entity.ID).
			Str("Username", entity.Username).
			Msg("access token verified")

		return entity, nil
	}

	return nil, errors.Wrap(InvalidTokenErr, "TokenService.VerifyRefreshToken")
}

func (s *TokenService) validateAccessSigningMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.Wrap(UnexpectedSigningMethodErr, "TokenService.validateAccessSigningMethod")
	}

	return []byte(s.cfg.AccessSecret), nil
}

func (s *TokenService) validateRefreshSigningMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.Wrap(UnexpectedSigningMethodErr, "TokenService.validateRefreshSigningMethod")
	}

	return []byte(s.cfg.RefreshSecret), nil
}
