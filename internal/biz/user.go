package biz

import (
	"context"
	"crypto/rand"
	"fmt"

	"kchat-be/internal/biz/entities"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUser(ctx context.Context, userID string) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	DeleteUser(ctx context.Context, userID string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*entities.User, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, username, email, password string) (*entities.User, error) {
	// Check if email exists
	existingUser, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Check if username exists
	existingUser, err = uc.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Generate user ID
	userID, err := generateUserID()
	if err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		UserID:   userID,
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		IsOnline: false,
	}

	// Validate
	if !user.ValidateEmail() {
		return nil, fmt.Errorf("invalid email format")
	}

	return uc.repo.CreateUser(ctx, user)
}

func (uc *UserUsecase) GetUser(ctx context.Context, userID string) (*entities.User, error) {
	return uc.repo.GetUser(ctx, userID)
}

func (uc *UserUsecase) UpdateUserStatus(ctx context.Context, userID string, online bool) (*entities.User, error) {
	user, err := uc.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	user.SetOnlineStatus(online)
	return uc.repo.UpdateUser(ctx, user)
}

func (uc *UserUsecase) ListUsers(ctx context.Context, page, pageSize int) ([]*entities.User, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return uc.repo.ListUsers(ctx, pageSize, offset)
}

func (uc *UserUsecase) AuthenticateUser(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Don't return password
	user.Password = ""
	return user, nil
}

func generateUserID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("user_%x", bytes), nil
}
