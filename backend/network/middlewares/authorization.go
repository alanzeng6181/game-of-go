package middlewares

import (
	"context"
	"net/http"

	security "github.com/alanzeng6181/game-of-go/security"
)

func ApplyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Token")
		if tokenStr == "" {
			tokenStr = r.URL.Query().Get("token")
		}
		if tokenStr == "" {
			w.WriteHeader(401)
			w.Write([]byte("missing auth token"))
		}
		userId, err := security.GetUserId(tokenStr)

		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte("unable to get userId from jwt token"))
			return
		}

		ctx := context.WithValue(r.Context(), "user", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
