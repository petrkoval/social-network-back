package middlewares

import (
	"github.com/petrkoval/social-network-back/internal/logger"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	l := logger.NewLogger()

	l.Debug().Msg("init logger middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Debug().Msgf("handling %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		l.Debug().Msgf("responsing %s %s", r.Method, r.RequestURI)
	})
}
