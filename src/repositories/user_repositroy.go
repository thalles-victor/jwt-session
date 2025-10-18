package repositories

import (
	"jwt-session/src/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetById(id string, user *models.User) error {
	err := r.db.Get(&user, `SELECT name, email, password, created_at FROM users WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string, user *models.User) error {
	err := r.db.Get(&user, `SELECT name, email, password, created_at FROM users WHERE email = ?`, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Create(user *models.User) error {
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO users (id, name, email, password, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.Password, user.CreatedAt)

	return err
}
