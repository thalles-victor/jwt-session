package repositories

import (
	"jwt-session/src/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db}
}

// GetAllByUserID busca todas as sessões de um usuário e preenche o slice passado como parâmetro
func (r *SessionRepository) GetAllByUserID(userID string, sessions *[]models.Session) error {
	return r.db.Select(sessions, `
		SELECT id, user_id, browser, ip, created_at, expires_at
		FROM sessions
		WHERE user_id = $1
	`, userID)
}

// GetByID busca uma sessão específica pelo ID e preenche o ponteiro passado como parâmetro
func (r *SessionRepository) GetByID(id string, session *models.Session) error {
	return r.db.Get(session, `
		SELECT id, user_id, browser, ip, created_at, expires_at
		FROM sessions
		WHERE id = $1
	`, id)
}

// Create insere uma nova sessão no banco
func (r *SessionRepository) Create(session *models.Session) (*models.Session, error) {
	// Se CreatedAt não estiver definido, define como agora
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO sessions (id, user_id, browser, ip, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, browser, ip, created_at, expires_at
	`

	var createdSession models.Session
	err := r.db.Get(&createdSession, query,
		session.ID,
		session.UserID,
		session.Browser,
		session.IP,
		session.CreatedAt,
		session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &createdSession, nil
}

// DeleteByID deleta uma sessão pelo ID
func (r *SessionRepository) DeleteByID(id string) error {
	_, err := r.db.Exec(`
		DELETE FROM sessions
		WHERE id = $1
	`, id)
	return err
}

// DeleteByID deleta uma sessão pelo ID
func (r *SessionRepository) DeleteByUserId(userId string) error {
	_, err := r.db.Exec(`
		DELETE FROM sessions
		WHERE user_id = $1
	`, userId)
	return err
}
