package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/geraldfingburke/dossier/server/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

// Service handles authentication
type Service struct {
	jwtSecret []byte
}

// NewService creates a new auth service
func NewService(jwtSecret string) *Service {
	if jwtSecret == "" {
		jwtSecret = "development-secret-key-change-in-production"
	}
	return &Service{
		jwtSecret: []byte(jwtSecret),
	}
}

// Register creates a new user
func (s *Service) Register(ctx context.Context, db *sql.DB, email, password, name string) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Create user
	var user models.User
	err = db.QueryRowContext(ctx, `
		INSERT INTO users (email, password, name)
		VALUES ($1, $2, $3)
		RETURNING id, email, name, created_at, updated_at
	`, email, string(hashedPassword), name).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &user, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, db *sql.DB, email, password string) (string, *models.User, error) {
	var user models.User
	var hashedPassword string

	err := db.QueryRowContext(ctx, `
		SELECT id, email, password, name, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &hashedPassword, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, fmt.Errorf("error finding user: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("error generating token: %w", err)
	}

	return tokenString, &user, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *Service) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

// GetUserFromContext retrieves user ID from context
func GetUserFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("user_id").(int)
	return userID, ok
}
