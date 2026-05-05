package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
)

func AuthMiddleware(logger *slog.Logger, secretKey []byte) func(http.Handler) http.Handler {
	logger = logger.With("module", "auth-middleware")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/auth/") { // todo: os okay? or can be better
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("missing authorization header", slog.String("path", r.URL.Path))
				writeError(w, application.ErrInvalidCredentials)
				return
			}

			tokenStr, ok := strings.CutPrefix(authHeader, "Bearer ")
			if !ok {
				logger.Warn("invalid authorization header format", slog.String("path", r.URL.Path))
				writeError(w, application.ErrInvalidCredentials)
				return
			}

			claims := &jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return secretKey, nil
			})
			if err != nil || !token.Valid {
				logger.Warn("invalid token", slog.String("path", r.URL.Path), slog.Any("err", err))
				writeError(w, application.ErrInvalidCredentials)
				return
			}

			sub, ok := (*claims)["sub"].(string)
			if !ok {
				logger.Warn("missing sub claim in token", slog.String("path", r.URL.Path))
				writeError(w, application.ErrInvalidCredentials)
				return
			}

			userID, err := uuid.Parse(sub)
			if err != nil {
				logger.Warn("invalid sub claim in token", slog.String("path", r.URL.Path), slog.Any("err", err))
				writeError(w, application.ErrInvalidCredentials)
				return
			}

			ctx := context.WithValue(r.Context(), idKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
