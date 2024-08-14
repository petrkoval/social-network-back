package middlewares

import (
	"github.com/go-chi/cors"
	"net/http"
)

var CorsMiddleware = cors.Handler(cors.Options{
	AllowedOrigins: []string{
		"http://localhost:5173",
	},
	AllowedMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	},
	AllowedHeaders: []string{
		"Authorization",
		"Content-Type",
	},
	ExposedHeaders: []string{
		"Content-Type",
	},
	AllowCredentials: true,
})
