package auth

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Service interface {
	Register(dto RegisterUserDTO) (Response, error)
	Login(dto LoginUserDTO) (Response, error)
	Logout(refreshToken string) error
	Refresh(refreshToken string) (Response, error)
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
