package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenService struct {
	accessSecret  []byte
	refreshSecret []byte
	privateKey    *rsa.PrivateKey
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`

	// Internal metadata used to persist/rotate the refresh token.
	RefreshJTI       string    `json:"-"`
	RefreshExpiresAt time.Time `json:"-"`
}

type Claims struct {
	jwt.RegisteredClaims
	TokenType string `json:"token_type,omitempty"`
}

func NewTokenService() (*TokenService, error) {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if accessSecret == "" {
		return nil, errors.New("JWT_ACCESS_SECRET not set")
	}

	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = accessSecret + "_refresh" // Fallback
	}

	// Generate RSA key for URL signing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	return &TokenService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		privateKey:    privateKey,
	}, nil
}

func (ts *TokenService) GenerateTokenPair(userID string) (*TokenPair, error) {
	// Access token (15 minutes)
	accessClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "objexrg",
			Subject:   userID,
		},
		TokenType: "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(ts.accessSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token (7 days)
	refreshJTI := uuid.NewString()
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "objexrg",
			Subject:   userID,
			ID:        refreshJTI,
		},
		TokenType: "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(ts.refreshSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:      accessTokenString,
		RefreshToken:     refreshTokenString,
		TokenType:        "Bearer",
		ExpiresIn:        15 * 60, // 15 minutes in seconds
		RefreshJTI:       refreshJTI,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

func (ts *TokenService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return ts.accessSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.TokenType != "access" {
			return nil, errors.New("not an access token")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (ts *TokenService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return ts.refreshSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.TokenType != "refresh" {
			return nil, errors.New("not a refresh token")
		}
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}

func (ts *TokenService) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
	claims, err := ts.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return ts.GenerateTokenPair(claims.Subject)
}

func (ts *TokenService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (ts *TokenService) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (ts *TokenService) SignURL(url string, expires time.Time) (string, error) {
	// Create a simple signature using SHA256
	h := sha256.New()
	h.Write([]byte(url + expires.Format(time.RFC3339) + string(ts.accessSecret)))
	signature := hex.EncodeToString(h.Sum(nil))

	signedURL := fmt.Sprintf("%s&signature=%s&expires=%d", url, signature, expires.Unix())
	return signedURL, nil
}

func (ts *TokenService) VerifySignedURL(url string, signature string, expires int64) bool {
	if time.Now().Unix() > expires {
		return false
	}

	h := sha256.New()
	h.Write([]byte(url + time.Unix(expires, 0).Format(time.RFC3339) + string(ts.accessSecret)))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return signature == expectedSignature
}

func GetUser(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", false
	}
	return userID.(string), true
}
