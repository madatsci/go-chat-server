package repositories

import (
	"github.com/go-pg/pg"
	"github.com/madatsci/go-chat-server/internal/models"
)

type (
	// User describes methods for storing user entity
	User interface {
		Create(email, password string) (*models.User, error)
		FindByEmail(email string) (*models.User, error)
	}

	userRepository struct {
		db *pg.DB
	}
)

// NewUserRepository creates new user repository instance
func NewUserRepository(db *pg.DB) User {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(email, hashedPassword string) (*models.User, error) {
	user := models.User{
		Email:    email,
		Password: hashedPassword,
	}

	if _, err := u.db.Model(&user).Where("email=?", email).SelectOrInsert(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) findBy(condition string, val interface{}) (*models.User, error) {
	var user models.User

	if err := u.db.Model(&user).Where(condition, val).First(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) FindByEmail(email string) (*models.User, error) {
	return u.findBy("email=?", email)
}
