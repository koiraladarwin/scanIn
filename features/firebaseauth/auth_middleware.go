package firebaseauth

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

type contextKey string

const (
	FirebaseUserContextKey contextKey = "firebaseUser"
	AccessLevelContextKey  contextKey = "accessLevel"
)

func FbUserFromContext(ctx context.Context) *auth.UserRecord {
	user, _ := ctx.Value(FirebaseUserContextKey).(*auth.UserRecord)
	return user
}

func AccessLevelFromContext(ctx context.Context) int {
	level, ok := ctx.Value(AccessLevelContextKey).(int)
	if !ok {
		return 0
	}
	return level
}

var accessLevels = map[string]int{
	"darwinkoirala123@gmail.com":       10,
	"darwinmage98422@gmail.com":        1,
	"darwinisabot7@gmail.com":        2,
	"ocsbusinesssolution@gmail.com":    2,
	"chhetrinirmal765@gmail.com":       2,
	"scan1@ocsbusinesssolution.com.np": 1,
	"scan2@ocsbusinesssolution.com.np": 1,
}

func (f *FirebaseAuth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "Malformed Authorization header", http.StatusUnauthorized)
			return
		}

		user, err := f.GetUserInfoByIDToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		level, allowed := accessLevels[user.Email]
		if !allowed {
			http.Error(w, "Unauthorized email", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), FirebaseUserContextKey, user)
		ctx = context.WithValue(ctx, AccessLevelContextKey, level)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
