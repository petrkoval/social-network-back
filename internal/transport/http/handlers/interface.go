package handlers

import (
	"github.com/petrkoval/social-network-back/internal/transport/http"
)

type Handler interface {
	MountOn(router *http.Router)
}
