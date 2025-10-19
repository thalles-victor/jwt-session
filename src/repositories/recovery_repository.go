package repositories

import (
	"jwt-session/src/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type RecoveryRepository struct {
	db *sqlx.DB
}

func NewRecoveryRepository(db *sqlx.DB) *RecoveryRepository {
	return &RecoveryRepository{db}
}

// GetByUserID busca um recovery pelo ID do usuário
func (r *RecoveryRepository) GetByUserID(userID string, recovery *models.Recovery) error {
	return r.db.Get(recovery, `
		SELECT id, user_id, email, code, attempts, expires_at, created_at, updated_at, expired
		FROM recoveries
		WHERE user_id = $1
	`, userID)
}

// GetByEmail busca um recovery pelo email
func (r *RecoveryRepository) GetByEmail(email string, recovery *models.Recovery) error {
	return r.db.Get(recovery, `
		SELECT id, user_id, email, code, attempts, expires_at, created_at, updated_at, expired
		FROM recoveries
		WHERE email = $1
	`, email)
}

// GetByCode busca um recovery pelo código
func (r *RecoveryRepository) GetByCode(code string, recovery *models.Recovery) error {
	return r.db.Get(recovery, `
		SELECT id, user_id, email, code, attempts, expires_at, created_at, updated_at, expired
		FROM recoveries
		WHERE code = $1
	`, code)
}

// Create insere um novo registro de recovery
func (r *RecoveryRepository) Create(recovery *models.Recovery) (models.Recovery, error) {
	if recovery.CreatedAt.IsZero() {
		recovery.CreatedAt = time.Now()
	}
	recovery.UpdatedAt = time.Now()

	query := `
		INSERT INTO recoveries (user_id, email, code, attempts, expires_at, created_at, updated_at, expired)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, email, code, attempts, expires_at, created_at, updated_at, expired
	`

	var createdRecovery models.Recovery
	if err := r.db.Get(&createdRecovery, query,
		recovery.UserID,
		recovery.Email,
		recovery.Code,
		recovery.Attempts,
		recovery.ExpiresAt,
		recovery.CreatedAt,
		recovery.UpdatedAt,
		recovery.Expired,
	); err != nil {
		return models.Recovery{}, err
	}

	return createdRecovery, nil
}

// IncrementAttempts aumenta o contador de tentativas
func (r *RecoveryRepository) IncrementAttempts(id int) error {
	_, err := r.db.Exec(`
		UPDATE recoveries
		SET attempts = attempts + 1, updated_at = $2
		WHERE id = $1
	`, id, time.Now())
	return err
}

// ClearByID limpa um recovery específico pelo id
func (r *RecoveryRepository) ClearByID(id int) error {
	_, err := r.db.Exec(`
		UPDATE recoveries
		SET code = NULL,
			attempts = 0,
			expires_at = NULL,
			updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

// ClearByUserID limpa um recovery específico pelo user_id
func (r *RecoveryRepository) ClearByUserID(userID string) error {
	_, err := r.db.Exec(`
		UPDATE recoveries
		SET code = NULL,
			attempts = 0,
			expires_at = NULL,
			updated_at = NOW()
		WHERE user_id = $1
	`, userID)
	return err
}

// ClearByEmail limpa um recovery específico pelo email
func (r *RecoveryRepository) ClearByEmail(email string) error {
	_, err := r.db.Exec(`
		UPDATE recoveries
		SET code = NULL,
			attempts = 0,
			expires_at = NULL,
			updated_at = NOW()
		WHERE email = $1
	`, email)
	return err
}

// UpdateRecovery atualiza code, attempts, expires_at e updated_at de um recovery existente
func (r *RecoveryRepository) UpdateRecovery(id int, code string, attempts int, expiresAt time.Time, expired bool) error {
	_, err := r.db.Exec(`
		UPDATE recoveries
		SET code = $1,
			attempts = $2,
			expires_at = $3,
			expired = $4,
			updated_at = NOW()
		WHERE id = $5
	`, code, attempts, expiresAt, expired, id)

	return err
}

// MarkAsExpired marca um recovery como expirado pelo id
func (r *RecoveryRepository) MarkAsExpired(id int) error {
	_, err := r.db.Exec(`
		UPDATE recoveries
		SET expired = TRUE,
			updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}
