package session

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (repo *repository) Save(sessID string, userID int64) error {
	if _, err := repo.db.Exec(`
	INSERT INTO sessions (id, user_id)
	VALUES ($1, $2);`, sessID, userID); err != nil {
		return fmt.Errorf("session.Save - error after processing the INSERT query: %w", err)
	}

	return nil
}

func (repo *repository) FindByID(id string) (*Session, error) {
	var s Session
	if err := repo.db.Get(&s, `
	SELECT id, user_id FROM sessions
	WHERE id = $1;`, id); err != nil {
		return nil, fmt.Errorf("session.FindByID - error after processing the SELECT query: %w", err)
	}

	return &s, nil
}

func (repo *repository) Delete(sessID string) error {
	if _, err := repo.db.Exec(`
	DELETE FROM sessions WHERE id = $1;`, sessID); err != nil {
		return fmt.Errorf("session.Delete - error after processing the DELETE query: %w", err)
	}

	return nil
}
