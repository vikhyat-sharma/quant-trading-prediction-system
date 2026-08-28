package util

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	JWTSecretEnvKey = "JWT_SECRET"
	TokenExpiry     = time.Hour * 24 // 24 hours; use refresh tokens for longer sessions
)

// GetJWTSecret returns the JWT secret from the environment.
// Returns an empty string if not configured; callers must treat this as a
// configuration error in production.
func GetJWTSecret() string {
	return strings.TrimSpace(os.Getenv(JWTSecretEnvKey))
}

// Claims represents JWT claims
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// GenerateJWT generates a signed JWT token.
// Returns an error if JWT_SECRET is not configured.
func GenerateJWT(userID int, email, role string) (string, error) {
	secret := GetJWTSecret()
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable is not configured")
	}
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// VerifyJWT verifies and parses a JWT token.
// Returns an error if JWT_SECRET is not configured or the token is invalid.
func VerifyJWT(tokenString string) (*Claims, error) {
	secret := GetJWTSecret()
	if secret == "" {
		return nil, errors.New("JWT_SECRET environment variable is not configured")
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// GenerateAPIKey generates a random API key
func GenerateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// IsAdminRole checks if role has admin privileges
func IsAdminRole(role string) bool {
	return role == "admin" || role == "superadmin"
}

// IsUserRole checks if role is a regular user
func IsUserRole(role string) bool {
	return role == "user"
}
