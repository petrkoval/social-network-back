package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/petrkoval/social-network-back/internal/domain"
	http2 "github.com/petrkoval/social-network-back/internal/transport/http"
	"github.com/petrkoval/social-network-back/internal/transport/http/middlewares"
	"github.com/rs/zerolog"
	"net/http"
)

const (
	path       = "/channels"
	channelUrl = "/{id}"
)

type ChannelService interface {
	FindAll() (*[]domain.Channel, error)
	FindByUserID(userID string) (*[]domain.Channel, error)
	FindByID(id string) (*domain.Channel, error)
	Create(dto domain.CreateChannelDTO) (*domain.Channel, error)
	Update(id string, dto domain.UpdateChannelDTO) (*domain.Channel, error)
	Delete(id string) error
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
		r.With(middlewares.Auth).Put("/", h.Update)
		r.With(middlewares.Auth).Delete("/", h.Delete)
	})

	router.Mount(path, h.router)
}

func (h *channelHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	entities, err := h.service.FindAll()

	if err != nil {
		switch {
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

	entities, err := h.service.FindByUserID(userID)

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

	entity, err := h.service.FindByID(id)

	if err != nil {
		switch {
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

	entity, err := h.service.Create(dto)

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

	entity, err := h.service.Update(id, dto)

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

	err := h.service.Delete(id)

	if err != nil {
		switch {
		default:
			WriteErrorResponse(w, r, err, http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
