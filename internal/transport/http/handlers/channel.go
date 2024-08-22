package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/petrkoval/social-network-back/internal/services"
	"github.com/petrkoval/social-network-back/internal/storage"
	http2 "github.com/petrkoval/social-network-back/internal/transport/http"
	"github.com/petrkoval/social-network-back/internal/transport/http/middlewares"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net/http"
)

const (
	path       = "/channels"
	channelUrl = "/{id}"
)

type ChannelService interface {
	FindAll(ctx context.Context, limit, offset string) ([]*domain.Channel, error)
	FindByUserID(ctx context.Context, userID string) ([]*domain.Channel, error)
	FindByID(ctx context.Context, id string) (*domain.Channel, error)
	Create(ctx context.Context, dto domain.CreateChannelDTO) (*domain.Channel, error)
	Update(ctx context.Context, id string, dto domain.UpdateChannelDTO) (*domain.Channel, error)
	Delete(ctx context.Context, id string) error
}

type channelHandler struct {
	service ChannelService
	logger  *zerolog.Logger
	router  *chi.Mux
}

func NewChannelHandler(s ChannelService, l *zerolog.Logger) Handler {
	r := chi.NewRouter()

	return &channelHandler{
		service: s,
		logger:  l,
		router:  r,
	}
}

func (h *channelHandler) MountOn(router *http2.Router) {
	h.router.Get("/", h.FindAll)

	h.router.Route(channelUrl, func(r chi.Router) {
		r.Get("/", h.FindByUserID)
		r.Get("/", h.FindByID)
		r.With(middlewares.Auth).Post("/", h.Create)
		r.With(middlewares.Auth).Patch("/", h.Update)
		r.With(middlewares.Auth).Delete("/", h.Delete)
	})

	router.Mount(path, h.router)
}

func (h *channelHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		limit  = chi.URLParam(r, "limit")
		offset = chi.URLParam(r, "offset")
	)

	entities, err := h.service.FindAll(r.Context(), limit, offset)

	if err != nil {
		switch {
		case errors.Is(err, services.QueryParamParsingErr):
			WriteErrorResponse(w, r, err, http.StatusBadRequest)
		default:
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(entities)
}

func (h *channelHandler) FindByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		userID = chi.URLParam(r, "id")
	)

	entities, err := h.service.FindByUserID(r.Context(), userID)

	if err != nil {
		switch {
		default:
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(entities)
}

func (h *channelHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		id = chi.URLParam(r, "id")
	)

	entity, err := h.service.FindByID(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, storage.NotFoundChannelErr):
			WriteErrorResponse(w, r, err, http.StatusNotFound)
		default:
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(entity)
}

func (h *channelHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		userID = chi.URLParam(r, "id")
		dto    = domain.CreateChannelDTO{UserID: userID}
		_      = json.NewDecoder(r.Body).Decode(&dto)
	)

	entity, err := h.service.Create(r.Context(), dto)

	if err != nil {
		switch {
		default:
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(entity)
}

func (h *channelHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		id  = chi.URLParam(r, "id")
		dto = domain.UpdateChannelDTO{}
		_   = json.NewDecoder(r.Body).Decode(&dto)
	)

	entity, err := h.service.Update(r.Context(), id, dto)

	if err != nil {
		switch {
		default:
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(entity)
}

func (h *channelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		id = chi.URLParam(r, "id")
	)

	err := h.service.Delete(r.Context(), id)

	if err != nil {
		switch {
		default:
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
