package middleware

import (
	"context"
	"net/http"
	"strings"

	"news-api/internal/models"
	"news-api/pkg/token"
	"news-api/utils"
)

type ctxKey string

const CtxActor ctxKey = "actor"

func AuthMiddleware(jwtManager *token.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				utils.WriteError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
				return
			}

			raw := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtManager.ValidateToken(raw)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			uidFloat, ok := claims["user_id"].(float64)
			if !ok {
				utils.WriteError(w, http.StatusUnauthorized, "invalid token payload")
				return
			}
			role, _ := claims["role"].(string)

			actor := models.Actor{UserID: int(uidFloat), Role: role}
			ctx := context.WithValue(r.Context(), CtxActor, actor)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
