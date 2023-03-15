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

// Generating token and send it to user's cookies
func AuthUser(w http.ResponseWriter, r *http.Request, userId int) error {
	token, err := generateToken(strconv.Itoa(userId))
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

// generateToken generates jwt token
func generateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Authorization": userId,
		"exp":           time.Now().Add(15 * time.Minute).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(signKey))
	if err != nil {

		err = fmt.Errorf("server error")
		return "", err
	}

	return tokenStr, nil
}

// validateToken validating token and returning jwt token or error
func validateToken(tokenString string) (*jwt.Token, error) {
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

// getUserId extracts userId from token
func getUserId(token *jwt.Token) (string, error) {
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

// getToken gets token from cookies
func getToken(r *http.Request) (string, error) {
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

// IsAuthenticated checks if user has token,if token is present - redirects to main page
func IsAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	token, err := getToken(r)
	if err != nil || token == "" {
		return false
	}
	return true
}

// GetUserIdFromCookies return's user's id from cookies
func GetUserIdFromCookies(r *http.Request) (string, error) {
	stringToken, err := getToken(r)
	if err != nil {
		return "", err
	}
	jwtToken, err := validateToken(stringToken)
	if err != nil {
		return "", err
	}
	userId, err := getUserId(jwtToken)
	if err != nil {
		return "", err
	}
	fmt.Println(userId)
	return userId, nil
}
