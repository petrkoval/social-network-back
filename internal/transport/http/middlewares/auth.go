package middlewares

import (
	"context"
	"encoding/json"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
)

type errorMessage struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
	URL        string `json:"url,omitempty"`
}

func writeErrorResponse(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorMessage{
		StatusCode: statusCode,
		Message:    err.Error(),
		URL:        r.Host + r.URL.Path,
	})
}

type service interface {
	VerifyAccessToken(accessToken string) (*domain.AuthUser, error)
}

func Auth(next http.Handler, s service, l *zerolog.Logger) http.Handler {

	l.Debug().Msg("init auth middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			writeErrorResponse(w, r, errors.New("authorization header is empty"), http.StatusUnauthorized)
			return
		}

		user, err := s.VerifyAccessToken(strings.Split(token, " ")[1])
		if err != nil {
			writeErrorResponse(w, r, errors.New("authorization header is empty"), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
