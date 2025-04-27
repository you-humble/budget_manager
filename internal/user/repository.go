package user

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

func (repo *repository) Save(u User) (int64, error) {
	var id int64
	if err := repo.db.Get(&id, `
	INSERT INTO users
	(login, password)
	VALUES
	($1, $2)
	RETURNING id;
	`, u.Login, u.Password); err != nil {
		return 0, fmt.Errorf("user.Save - error after processing the INSERT query: %w", err)
	}

	return id, nil
}

func (repo *repository) FindByLogin(login string) (User, error) {
	var u User
	if err := repo.db.Get(&u, `
	SELECT id, login, password FROM users
	WHERE login = $1;`, login); err != nil {
		return User{}, fmt.Errorf("user.FindByLogin - error after processing the SELECT query: %w", err)
	}

	return u, nil
}
