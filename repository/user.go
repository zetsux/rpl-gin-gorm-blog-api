package repository

import (
	"context"
	"errors"
	"go-blogrpl/entity"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	// db transaction
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(ctx context.Context, tx *gorm.DB) error
	RollbackTx(ctx context.Context, tx *gorm.DB)

	// functional
	CreateNewUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	FindByUsernameOrEmail(ctx context.Context, tx *gorm.DB, username string, email string) (entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (userR *userRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := userR.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (userR *userRepository) CommitTx(ctx context.Context, tx *gorm.DB) error {
	err := tx.WithContext(ctx).Commit().Error
	if err == nil {
		return err
	}
	return nil
}

func (userR *userRepository) RollbackTx(ctx context.Context, tx *gorm.DB) {
	tx.WithContext(ctx).Debug().Rollback()
}

func (userR *userRepository) CreateNewUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	var err error
	if tx == nil {
		tx = userR.db.WithContext(ctx).Debug().Create(&user)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Create(&user).Error
	}

	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (userR *userRepository) FindByUsernameOrEmail(ctx context.Context, tx *gorm.DB, username string, email string) (entity.User, error) {
	var err error
	var user entity.User
	if tx == nil {
		tx = userR.db.WithContext(ctx).Debug().Where("username = $1 OR email = $2", username, email).Take(&user)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("username = $1 OR email = $2", username, email).Take(&user).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return user, err
	}
	return user, nil
}

func (userR *userRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User

	tx := userR.db.WithContext(ctx).Debug().Preload("Blogs").Preload("BlogLikes").Preload("CommentLikes").Find(&users)
	if tx.Error != nil && !(errors.Is(tx.Error, gorm.ErrRecordNotFound)) {
		return users, tx.Error
	}
	return users, nil
}

func (userR *userRepository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User

	tx := userR.db.WithContext(ctx).Debug().Where("username = $1", username).Preload("Blogs").Preload("BlogLikes").Preload("CommentLikes").Take(&user)
	if tx.Error != nil && !(errors.Is(tx.Error, gorm.ErrRecordNotFound)) {
		return user, tx.Error
	}
	return user, nil
}
