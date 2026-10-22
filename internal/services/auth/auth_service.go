package auth

import (
	"errors"
	"fmt"
	dtos "objexrg/internal/dtos/auth/input"
	"time"

	"objexrg/internal/interfaces"
	models "objexrg/internal/models/auth"
	bucketsvc "objexrg/internal/services/bucket"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	auth         interfaces.AuthInterface
	db           *gorm.DB
	tokenService *TokenService
	bucketSvc    *bucketsvc.BucketService
}

func NewAuthService(auth interfaces.AuthInterface, db *gorm.DB, tokenService *TokenService, bucketSvc *bucketsvc.BucketService) *AuthService {
	return &AuthService{auth: auth, db: db, tokenService: tokenService, bucketSvc: bucketSvc}
}

func (s *AuthService) Register(req *dtos.RegisterRequest) error {
	// Parsear birth date
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return errors.New("invalid birth date format")
	}

	hashedPassword, err := s.tokenService.HashPassword(req.Password)
	if err != nil {
		return err
	}

	userID, err := s.auth.RegisterUser(req.Email, req.Username, string(hashedPassword), req.FirstName, req.LastName, birthDate)
	if err != nil {
		return err
	}

	// Crear bucket raíz para el usuario
	if err := s.bucketSvc.CreateRootBucket(userID); err != nil {
		return fmt.Errorf("failed to create root bucket: %w", err)
	}

	return nil
}

func (s *AuthService) Login(req *dtos.LoginRequest) (*models.UserCredentials, *TokenPair, error) {
	user, err := s.auth.GetByEmail(req.Email)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Verificar contraseña usando el token service
	if !s.tokenService.CheckPassword(user.Password, req.Password) {
		return nil, nil, errors.New("invalid credentials")
	}

	// Generar tokens
	tokens, err := s.tokenService.GenerateTokenPair(user.ID.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	if err := s.persistRefreshToken(user.ID, tokens); err != nil {
		return nil, nil, fmt.Errorf("failed to persist refresh token: %w", err)
	}

	return user, tokens, nil
}

func (s *AuthService) Create(user *models.UserCredentials) error {
	return s.auth.Create(user)
}

func (s *AuthService) GetByEmail(email string) (*models.UserCredentials, error) {
	return s.auth.GetByEmail(email)
}

func (s *AuthService) GetByUsername(username string) (*models.UserCredentials, error) {
	return s.auth.GetByUsername(username)
}

func (s *AuthService) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// The refresh token must still exist and not have been revoked/rotated.
	stored, err := s.auth.GetRefreshTokenByID(claims.ID)
	if err != nil {
		return nil, err
	}
	if stored == nil || stored.Revoked {
		return nil, errors.New("refresh token is no longer valid")
	}
	if stored.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	// Rotate: revoke the used token before issuing a new pair.
	if err := s.auth.RevokeRefreshToken(claims.ID); err != nil {
		return nil, fmt.Errorf("failed to rotate refresh token: %w", err)
	}

	tokens, err := s.tokenService.GenerateTokenPair(claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("invalid subject in refresh token: %w", err)
	}

	if err := s.persistRefreshToken(userID, tokens); err != nil {
		return nil, fmt.Errorf("failed to persist refresh token: %w", err)
	}

	return tokens, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		// Already invalid/expired: treat as a successful logout.
		return nil
	}
	return s.auth.RevokeRefreshToken(claims.ID)
}

func (s *AuthService) persistRefreshToken(userID uuid.UUID, tokens *TokenPair) error {
	return s.auth.CreateRefreshToken(&models.RefreshToken{
		ID:        uuid.MustParse(tokens.RefreshJTI),
		UserID:    userID,
		ExpiresAt: tokens.RefreshExpiresAt,
		Revoked:   false,
		CreatedAt: time.Now(),
	})
}
