package logging

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// Logging is the function that wraps a http.Handler to provide the log information from the incoming request.
func Logging() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			handlers.LoggingHandler(os.Stdout, next).ServeHTTP(w, r)
		}
	}
}
