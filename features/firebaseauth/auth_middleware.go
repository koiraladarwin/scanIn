package firebaseauth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

type contextKey string

const (
	FirebaseUserContextKey contextKey = "firebaseUser"
)

func FbUserFromContext(ctx context.Context) (*auth.UserRecord,bool) {
	user, ok := ctx.Value(FirebaseUserContextKey).(*auth.UserRecord)
	return user,ok
}


func (f *FirebaseAuth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

		user, err := f.GetUserInfoByIDToken(r.Context(), token)

		if err != nil {
			log.Print("Firebase Authorization header: " + err.Error())
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
    
		ctx := context.WithValue(r.Context(), FirebaseUserContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
