package middleware

import (
	"net/http"

	"github.com/koiraladarwin/scanin/features/firebaseauth"
)

func RequireAccessLevel(minLevel int, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(firebaseauth.AccessLevelContextKey).(int)
		if !ok  || user < minLevel  {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
  
  	
		next(w, r)
	}
}

