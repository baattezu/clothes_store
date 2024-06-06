package auth

import (
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler, authClient *AuthClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		valid, userID, scope, err := authClient.ValidateToken(token)
		if err != nil || !valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем userID и scope в контекст запроса
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "scope", scope)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
