package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	signKey = os.Getenv("SigningKey")
)

// AuthMiddleWare
func AuthMiddleWare(next func() http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve token from request
		token, err := GetToken(r)
		if err != nil || token == "" {
			// redirecting to login page
			r.Method = "GET"
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next().ServeHTTP(w, r)
	})
}

// Generating token and send it to user in cookies
func AuthUser(w http.ResponseWriter, r *http.Request, userId int) error {
	token, err := GenerateToken(userId)
	if err != nil {
		return err
	}
	cookies := &http.Cookie{}
	cookies.Name = "Authorization"
	cookies.Value = "Bearer " + token
	cookies.Path = "/"
	cookies.Expires = time.Now().Add(15 * time.Minute)
	http.SetCookie(w, cookies)
	return nil
}

// GenerateToken generates jwt token
func GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Authorization": strconv.Itoa(userId),
		"exp":           time.Now().Add(15 * time.Minute).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(signKey))
	if err != nil {

		err = fmt.Errorf("server error")
		return "", err
	}

	return tokenStr, nil
}

// ValidateToken validating token and returning user's id or error
func ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the token's signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("bad signing method: %v", token.Header["alg"])
		}
		// Return the secret key used to sign the token
		return []byte(signKey), nil
	})
	if err != nil {
		return token, fmt.Errorf("failed to parse token: %v", err)
	}
	return token, nil
}

// GetUserId extracts userId from token
func GetUserId(token *jwt.Token) (string, error) {
	// Extract the userID field from the token's payload
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("unexpected claims format")
	}
	userId, ok := claims["Authorization"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid userID field")
	}
	return userId, nil
}

// GetToken gets token from cookies
func GetToken(r *http.Request) (string, error) {
	token, err := r.Cookie("Authorization")
	if err != nil {
		return "", err
	}
	splitToken := strings.Split(token.Value, " ")
	if len(splitToken) != 2 || splitToken[1] == "" {
		err := fmt.Errorf("Token not found, token: %v", token)
		return "", err
	}

	return splitToken[1], nil
}

// IsAuthenticated checks if user has token,if token is present - redirects to /chat/ page else
// redirects to login page
func IsAuthenticated(w http.ResponseWriter, r *http.Request) {
	token, err := GetToken(r)
	if err != nil || token == "" {
		return
	}
	http.Redirect(w, r, "/chat/", http.StatusFound)
}
