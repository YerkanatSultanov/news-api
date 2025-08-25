package middleware

import (
	"net/http"

	"news-api/pkg/logger"
	"news-api/utils"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Log.Error("panic recovered",
					"error", rec,
					"path", r.URL.Path,
				)
				utils.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
