package middleware

import (
	"context"
	"fmt"
	"net/http"
	"notes-api-go/utils"
	"slices"
)

type UserID struct{}
type Roles struct{}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenStr := authHeader[len("Bearer "):]

		claims, err := utils.VerifyToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user info to the request
		// r.Header.Set("UserID", strconv.FormatInt(claims.UserID, 10))
		// r.Header.Set("Role", claims.Role)

		// next(w, r)

		// ctx := context.WithValue(r.Context(), UserID{}, claims.UserID)
		// ctx = context.WithValue(ctx, Roles{}, claims.Roles)

		userID := claims.UserInfo.UserID
		roles := claims.UserInfo.Roles
		ctx := context.WithValue(r.Context(), UserID{}, userID)
		ctx = context.WithValue(ctx, Roles{}, roles)
		next(w, r.WithContext(ctx))
	}
}

func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// role := r.Header.Get("Role")
		roles, err := GetUserRoleFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !slices.Contains(roles, "admin") {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(UserID{}).(int64)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func GetUserRoleFromContext(ctx context.Context) ([]string, error) {
	var emptyResponse []string
	roles, ok := ctx.Value(Roles{}).([]string)
	if !ok {
		return emptyResponse, fmt.Errorf("user role not found in context")
	}
	return roles, nil
}
