package middlewares

import (
	"context"
	"errors"
	"github.com/petrkoval/social-network-back/internal/config"
	"github.com/petrkoval/social-network-back/internal/logger"
	"github.com/petrkoval/social-network-back/internal/services"
	"github.com/petrkoval/social-network-back/internal/transport/http/handlers"
	"net/http"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	l := logger.NewLogger()
	cfg, err := config.MustLoad()

	if err != nil {
		l.Panic().Err(err).Msg("error loading config in auth middleware")
	}

	s := services.NewTokenService(nil, l, cfg.Tokens)

	l.Debug().Msg("init auth middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			handlers.WriteErrorResponse(w, r, errors.New("authorization header is empty"), http.StatusUnauthorized)
		}

		user, err := s.VerifyAccessToken(strings.Split(token, " ")[1])
		if err != nil {
			handlers.WriteErrorResponse(w, r, err, http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
