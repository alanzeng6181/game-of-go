package middlewares

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
)

//TODO secure this
var JWTSigningKey []byte = []byte("TemporaryNotSecureKey")

func ApplyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Token")
		if tokenStr == "" {
			tokenStr = r.URL.Query().Get("Token")
		}
		if tokenStr == "" {
			w.WriteHeader(401)
			w.Write([]byte("missing auth token"))
		}
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("expected HS356 signing method")
			}
			return JWTSigningKey, nil
		})

		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte("unable to parse jwt token"))
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			context := context.WithValue(r.Context(), "user", claims["user"].(string))
			next.ServeHTTP(w, r.Context())
		} else {
			w.WriteHeader(401)
			w.Write([]byte("invalid token claims"))
		}

	})
}
