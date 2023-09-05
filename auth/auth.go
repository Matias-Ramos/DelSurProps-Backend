package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)



func ValidateJwt(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if clientToken := r.Header.Get("Authentication"); clientToken != "" {
				token, err := parseJwtToken(clientToken, jwtSecret, w)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized: " + err.Error()))
					return
				}
				if token.Valid {

					next.ServeHTTP(w, r)
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized"))
			}
		})
	}
}
/* 
 parseJwtToken parses a client-supplied token string into an actual Jwt,
 validates it using a HMAC-based signing method,
 and returns the parsed token or an error
*/
func parseJwtToken(clientToken string, jwtSecret string, w http.ResponseWriter) (*jwt.Token, error) {
	token, err := jwt.Parse(clientToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
			return nil, fmt.Errorf("not authorized")
		}
		return []byte(jwtSecret), nil
	})
	return token, err
}


func CreateJwt(jwtSecret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("err at CreateJwt(): %v", err.Error())
		return "", err
	}
	return tokenStr, nil
}
func GetJwt(apiKey, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// io.ReadCloser to []byte
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error on io.ReadAll() - GetJwt(): %s", err)
			return
		}

		// []byte to map
		var m map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &m); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error on json.Unmarshal() - []byte to map - GetJwt(): %s", err)
			return
		}
		frontEndKey := m["Access"]

		if frontEndKey != "" && frontEndKey == apiKey {
			token, err := CreateJwt(jwtSecret)
			if err != nil {
				log.Printf("token creation failed on auth.go - GetJwt(): %v", err)
				return
			}
			w.Write([]byte(token))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Wrong credentials - not authorized"))
		}
	}
}

