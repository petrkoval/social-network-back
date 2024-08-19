package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/petrkoval/social-network-back/internal/services"
	"github.com/petrkoval/social-network-back/internal/storage"
	http2 "github.com/petrkoval/social-network-back/internal/transport/http"
	"github.com/rs/zerolog"
	"net/http"
)

const (
	registerUrl = "/register"
	loginUrl    = "/login"
	logoutUrl   = "/logout"
	refreshUrl  = "/refresh"
)

type AuthService interface {
	Register(ctx context.Context, dto domain.CreateUserDTO) (*services.AuthResponse, error)
	Login(ctx context.Context, dto domain.CreateUserDTO) (*services.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (*services.AuthResponse, error)
}

type authHandler struct {
	service AuthService
	logger  *zerolog.Logger
	router  *chi.Mux
}

func NewAuthHandler(s AuthService, l *zerolog.Logger) Handler {
	r := chi.NewRouter()

	return &authHandler{
		service: s,
		logger:  l,
		router:  r,
	}
}

func (h *authHandler) MountOn(router *http2.Router) {
	h.router.Post(registerUrl, h.Register)
	h.router.Post(loginUrl, h.Login)
	h.router.Post(logoutUrl, h.Logout)
	h.router.Get(refreshUrl, h.Refresh)

	router.Mount("/", h.router)
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		entity domain.CreateUserDTO
	)

	_ = json.NewDecoder(r.Body).Decode(&entity)

	response, err := h.service.Register(r.Context(), entity)
	if err != nil {
		switch {
		case errors.Is(err, services.UserExistsErr):
			h.logger.Error().Stack().Err(err).Msg("User already exists")
			WriteErrorResponse(w, r, err, http.StatusConflict)
			return
		default:
			h.logger.Error().Stack().Err(err).Msg("unhandled error")
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
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

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		entity domain.CreateUserDTO
	)

	_ = json.NewDecoder(r.Body).Decode(&entity)

	response, err := h.service.Login(r.Context(), entity)
	if err != nil {
		switch {
		case errors.Is(err, storage.NotFoundUserErr):
			h.logger.Error().Stack().Err(err).Msg("User not found")
			WriteErrorResponse(w, r, err, http.StatusNotFound)
			return
		case errors.Is(err, services.WrongPasswordErr):
			WriteErrorResponse(w, r, err, http.StatusNotFound)
			return
		default:
			h.logger.Error().Stack().Err(err).Msg("unhandled error")
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
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

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			h.logger.Error().Stack().Err(err).Msg("no refresh_token cookie found")
			WriteErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		default:
			h.logger.Error().Stack().Err(err).Msg("unhandled error")
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	err = h.service.Logout(r.Context(), refreshToken.Value)
	if err != nil {
		h.logger.Error().Stack().Err(err).Msg("unhandled error")
		WriteErrorResponse(w, r, err, http.StatusInternalServerError)
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

func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			h.logger.Error().Stack().Err(err).Msg("no refresh_token cookie found")
			WriteErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		default:
			h.logger.Error().Stack().Err(err).Msg("unhandled error")
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	response, err := h.service.Refresh(r.Context(), refreshToken.Value)
	if err != nil {
		switch {
		case errors.Is(err, services.TokenExpiredErr):
			h.logger.Error().Stack().Err(err).Msg("refresh token expired")
			WriteErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		case errors.Is(err, services.InvalidTokenErr):
			h.logger.Error().Stack().Err(err).Msg("invalid token")
			WriteErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		case errors.Is(err, storage.NotFoundTokenErr):
			h.logger.Error().Stack().Err(err).Msg("token not found")
			WriteErrorResponse(w, r, err, http.StatusUnauthorized)
			return
		default:
			h.logger.Error().Stack().Err(err).Msg("unhandled error")
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
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
