package data

import (
	"context"
	"time"

	"kchat-be/internal/biz"
	"kchat-be/internal/biz/entities"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := r.data.db.WithContext(ctx).Create(user).Error; err != nil {
		r.log.Errorf("failed to create user: %v", err)
		return nil, err
	}

	// Don't return password
	user.Password = ""

	return user, nil
}

func (r *userRepo) GetUser(ctx context.Context, userID string) (*entities.User, error) {
	var user entities.User

	err := r.data.db.WithContext(ctx).
		Select("user_id, username, email, is_online, created_at, updated_at, profile_picture").
		Where("user_id = ?", userID).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		r.log.Errorf("failed to get user: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User

	err := r.data.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.Errorf("failed to get user by email: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User

	err := r.data.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.Errorf("failed to get user by username: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	user.UpdatedAt = time.Now()

	err := r.data.db.WithContext(ctx).
		Select("username", "email", "is_online", "updated_at", "profile_picture").
		Where("user_id = ?", user.UserID).
		Updates(user).Error

	if err != nil {
		r.log.Errorf("failed to update user: %v", err)
		return nil, err
	}

	// Get updated user
	return r.GetUser(ctx, user.UserID)
}

func (r *userRepo) DeleteUser(ctx context.Context, userID string) error {
	err := r.data.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&entities.User{}).Error

	if err != nil {
		r.log.Errorf("failed to delete user: %v", err)
		return err
	}

	return nil
}

func (r *userRepo) ListUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var users []*entities.User

	err := r.data.db.WithContext(ctx).
		Select("user_id, username, email, is_online, created_at, updated_at, profile_picture").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		r.log.Errorf("failed to list users: %v", err)
		return nil, err
	}

	return users, nil
}
