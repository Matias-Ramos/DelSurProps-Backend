package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func ValidateJwt(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if clientToken := r.Header.Get("Authentication"); clientToken != "" {
				token, err := parseJwtToken(clientToken, jwtSecret, w)
				if err != nil {
					fmt.Println("err1: ", err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized: " + err.Error()))
					return
				}
				if token.Valid {
					next.ServeHTTP(w, r)
				}
			} else {
				fmt.Println(`clientToken == "" `)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized"))
			}
		})
	}
}

/*
parseJwtToken parses a client-supplied token string into an actual Jwt,
validates it using a HMAC-based signing method,
and returns the parsed token or an error to ValidateJwt.
*/
func parseJwtToken(clientToken string, jwtSecret string, w http.ResponseWriter) (*jwt.Token, error) {
	token, err := jwt.Parse(clientToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Println(`t.Method !ok at parseJwtToken`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
			return nil, fmt.Errorf("not authorized")
		}
		return []byte(jwtSecret), nil
	})
	return token, err
}
