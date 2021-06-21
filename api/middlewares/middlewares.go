package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/haasin-farooq/todo-go-app/api/responses"
)

func SetContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func AuthJwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := map[string]interface{}{
			"status": "failed",
			"message": "Missing authorization token",
		}

		authHeader := r.Header.Get("Authorization")

		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		authHeader = strings.TrimSpace(authHeader)

		if authHeader == "" {
			responses.JSON(w, http.StatusBadRequest, res)
			return
		}

		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			res["status"] = "failed"
			res["message"] = "Invalid token, please login"
			responses.JSON(w, http.StatusUnauthorized, res)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		c := context.WithValue(r.Context(), "userID", claims["userID"])

		next.ServeHTTP(w, r.WithContext(c))
	})
}