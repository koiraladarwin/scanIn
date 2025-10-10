package firebaseauth

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"firebase.google.com/go/auth"
)

type contextKey string

const (
	FirebaseUserContextKey contextKey = "firebaseUser"
)

func FbUserFromContext(ctx context.Context) (*auth.UserRecord, bool) {
	user, ok := ctx.Value(FirebaseUserContextKey).(*auth.UserRecord)
	return user, ok
}

func (f *FirebaseAuth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mode := os.Getenv("MODE")

		if strings.ToLower(mode) == "dev" {
			mockUser := &auth.UserRecord{
				UserInfo: &auth.UserInfo{
					UID:         "mock-uid-12345",
					DisplayName: "Darwin Shrestha",
					Email:       "darwin@example.com",
					PhotoURL:    "https://example.com/images/darwin.jpg",
				},
				CustomClaims: map[string]interface{}{
					"role": "admin",
				},
			}

			ctx := context.WithValue(r.Context(), FirebaseUserContextKey, mockUser)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Print("Misssing Authorization header: " + authHeader)
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			log.Print("Malformed Authorization header: " + authHeader)
			http.Error(w, "Malformed Authorization header", http.StatusUnauthorized)
			return
		}

		user, err := f.VerifyUserByIdToken(r.Context(), token)

		if err != nil {
			log.Print("Firebase Authorization header: " + err.Error())
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

    ctx := context.WithValue(r.Context(), FirebaseUserContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
