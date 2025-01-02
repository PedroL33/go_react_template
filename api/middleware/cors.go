package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func (m middleWareManager) Cors(next http.Handler) http.Handler {
	// Set up CORS options (you can modify the allowed origins as needed)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	return corsHandler.Handler(next)
}
