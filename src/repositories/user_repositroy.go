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
	err := r.db.Get(user, `SELECT name, email, password, created_at FROM "users" WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string, user *models.User) error {
	err := r.db.Get(user, `SELECT id, name, email, password, created_at FROM "users" WHERE email = $1`, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Create(user *models.User) (models.User, error) {
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO "users" (id, name, email, password, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, email, password, created_at
	`

	var createdUser models.User
	if err := r.db.Get(&createdUser, query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
	); err != nil {
		return models.User{}, err
	}

	return createdUser, nil
}

// UpdatePassword atualiza a senha do usuário pelo ID
func (r *UserRepository) UpdatePassword(userID string, newPassword string) error {
	_, err := r.db.Exec(`
		UPDATE "users"
		SET password = $1, updated_at = $2
		WHERE id = $3
	`, newPassword, time.Now(), userID)

	return err
}
