package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/petrkoval/social-network-back/internal/services"
	"github.com/petrkoval/social-network-back/internal/storage"
	http2 "github.com/petrkoval/social-network-back/internal/transport/http"
	"github.com/petrkoval/social-network-back/internal/transport/http/handlers"
	"github.com/rs/zerolog"
	"net/http"
)

const (
	registerUrl = "/register"
	loginUrl    = "/login"
	logoutUrl   = "/logout"
	refreshUrl  = "/refresh"
)

type Service interface {
	Register(ctx context.Context, dto domain.CreateUserDTO) (*services.AuthResponse, error)
	Login(ctx context.Context, dto domain.CreateUserDTO) (*services.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (*services.AuthResponse, error)
}

type Handler struct {
	service Service
	logger  *zerolog.Logger
	router  *chi.Mux
}

func NewAuthHandler(s Service, l *zerolog.Logger) handlers.Handler {
	r := chi.NewRouter()

	return &Handler{
		service: s,
		logger:  l,
		router:  r,
	}
}

func (h *Handler) MountOn(router *http2.Router) {
	h.router.Post(registerUrl, h.Register)
	h.router.Post(loginUrl, h.Login)
	h.router.Post(logoutUrl, h.Logout)
	h.router.Get(refreshUrl, h.Refresh)

	router.Mount("/", h.router)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		entity domain.CreateUserDTO
	)

	_ = json.NewDecoder(r.Body).Decode(&entity)

	response, err := h.service.Register(r.Context(), entity)
	if err != nil {
		switch {
		case errors.Is(err, services.UserExistsErr):
			h.logger.Error().Stack().Err(err).Msg("User already exists")
			handlers.WriteErrorResponse(w, r, err, http.StatusConflict)
			return
		default:
			h.logger.Error().Stack().Err(err).Msg("unhandled error")
			handlers.WriteErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		entity domain.CreateUserDTO
	)

	_ = json.NewDecoder(r.Body).Decode(&entity)

	response, err := h.service.Login(r.Context(), entity)
	if err != nil {
		switch {
		case errors.Is(err, storage.NotFoundUserErr):
			h.logger.Error().Stack().Err(err).Msg("User not found")
			handlers.WriteErrorResponse(w, r, err, http.StatusNotFound)
			return
		case errors.Is(err, services.WrongPasswordErr):
			handlers.WriteErrorResponse(w, r, err, http.StatusNotFound)
			return
		default:
			h.logger.Error().Stack().Err(err).Msg("unhandled error")
			handlers.WriteErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			h.logger.Error().Stack().Err(err).Msg("no refresh_token cookie found")
			handlers.WriteErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		default:
			h.logger.Error().Stack().Err(err).Msg("unhandled error")
			handlers.WriteErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	err = h.service.Logout(r.Context(), refreshToken.Value)
	if err != nil {
		h.logger.Error().Stack().Err(err).Msg("unhandled error")
		handlers.WriteErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "refresh_token",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			h.logger.Error().Stack().Err(err).Msg("no refresh_token cookie found")
			handlers.WriteErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		default:
			h.logger.Error().Stack().Err(err).Msg("unhandled error")
			handlers.WriteErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	response, err := h.service.Refresh(r.Context(), refreshToken.Value)
	if err != nil {
		switch {
		case errors.Is(err, services.TokenExpiredErr):
			h.logger.Error().Stack().Err(err).Msg("refresh token expired")
			handlers.WriteErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		case errors.Is(err, services.InvalidTokenErr):
			h.logger.Error().Stack().Err(err).Msg("invalid token")
			handlers.WriteErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		case errors.Is(err, storage.NotFoundTokenErr):
			h.logger.Error().Stack().Err(err).Msg("token not found")
			handlers.WriteErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		default:
			h.logger.Error().Stack().Err(err).Msg("unhandled error")
			handlers.WriteErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
