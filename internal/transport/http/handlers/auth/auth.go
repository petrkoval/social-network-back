package auth

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/petrkoval/social-network-back/internal/services"
	"net/http"
)

type Service interface {
	Register(ctx context.Context, dto domain.CreateUserDTO) (*services.AuthResponse, error)
	Login(ctx context.Context, dto domain.CreateUserDTO) (*services.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (*services.AuthResponse, error)
}

type Handler struct {
	service Service
	Router  *chi.Mux
}

func NewAuthRouter(s Service) *Handler {
	r := chi.NewRouter()

	return &Handler{
		service: s,
		Router:  r,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {

}
