package handlers

import (
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
	Update(userID string, dto domain.UpdateChannelDTO) (*domain.Channel, error)
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
	h.router.Use(middlewares.Auth)

	h.router.Get("/", h.FindAll)
	h.router.Post("/", h.Create)

	h.router.Route(channelUrl, func(r chi.Router) {
		r.Get("/", h.FindByUserID)
		r.Get("/", h.FindByID)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})

	router.Mount(path, h.router)
}

func (h *channelHandler) FindAll(w http.ResponseWriter, r *http.Request) {

}

func (h *channelHandler) FindByUserID(w http.ResponseWriter, r *http.Request) {

}

func (h *channelHandler) FindByID(w http.ResponseWriter, r *http.Request) {

}

func (h *channelHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *channelHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h *channelHandler) Delete(w http.ResponseWriter, r *http.Request) {

}
