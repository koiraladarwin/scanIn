package firebaseauth

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

type contextKey string

const FirebaseUserContextKey = contextKey("firebaseUser")

func FbUserFromContext(ctx context.Context) *auth.UserRecord {
	user, _ := ctx.Value(FirebaseUserContextKey).(*auth.UserRecord)
	return user
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

		allowed := map[string]bool{
			"darwinkoirala123@gmail.com":    true,
			"ocsbusinesssolution@gmail.com": true,
			"chhetrinirmal765@gmail.com":    true,
		}

		if !allowed[user.Email] {
			http.Error(w, "Unauthorized email", http.StatusUnauthorized)
			return

		}

		ctx := context.WithValue(r.Context(), FirebaseUserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
