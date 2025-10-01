package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"time"

	"practice-go-crud/internal/domain"
	"practice-go-crud/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error)
	AdminRegister(ctx context.Context, req *domain.AdminRegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	ValidateToken(tokenString string) (*domain.User, error)
	Logout(ctx context.Context, tokenString string) error
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

func (s *authService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *authService) AdminRegister(ctx context.Context, req *domain.AdminRegisterRequest) (*domain.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateJWTToken(user)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &domain.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) generateJWTToken(user *domain.User) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET not configured")
	}

	expiryHours := os.Getenv("JWT_EXPIRY_HOURS")
	if expiryHours == "" {
		expiryHours = "24"
	}

	hours, err := strconv.Atoi(expiryHours)
	if err != nil {
		hours = 24
	}

	jti, err := s.generateJTI()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
		"jti":     jti,
		"exp":     time.Now().Add(time.Duration(hours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *authService) ValidateToken(tokenString string) (*domain.User, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		jti, ok := claims["jti"].(string)
		if !ok {
			return nil, errors.New("JTI not found in token")
		}

		isBlacklisted, err := s.tokenRepo.IsBlacklisted(context.Background(), jti)
		if err != nil {
			return nil, err
		}
		if isBlacklisted {
			return nil, ErrInvalidToken
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, ErrInvalidToken
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		name, ok := claims["name"].(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		role, ok := claims["role"].(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		return &domain.User{
			ID:    int(userID),
			Email: email,
			Name:  name,
			Role:  role,
		}, nil
	}

	return nil, ErrInvalidToken
}

func (s *authService) Logout(ctx context.Context, tokenString string) error {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return errors.New("JWT_SECRET not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		jti, ok := claims["jti"].(string)
		if !ok {
			return errors.New("JTI not found in token")
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return errors.New("user_id not found in token")
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return errors.New("exp not found in token")
		}

		expiredAt := time.Unix(int64(exp), 0)

		return s.tokenRepo.AddToBlacklist(ctx, jti, int(userID), expiredAt)
	}

	return ErrInvalidToken
}
