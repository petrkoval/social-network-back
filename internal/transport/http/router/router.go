package router

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/petrkoval/social-network-back/internal/config"
	"github.com/petrkoval/social-network-back/internal/transport/http/router/middlewares"
	"net/http"
	"time"
)

type Router struct {
	cfg *config.ServerConfig
	*chi.Mux
}

func NewRouter(cfg *config.ServerConfig) *Router {
	return &Router{
		cfg: cfg,
		Mux: chi.NewRouter(),
	}
}

func (r *Router) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", r.cfg.Port), r)
}

func (r *Router) InitMiddlewares() {
	r.Use(middlewares.Logger)
	r.Use(middlewares.CorsMiddleware)
	r.Use(middleware.Timeout(time.Duration(r.cfg.WriteTimeout)))
	r.Use(middleware.Recoverer)
}
