package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var schema string = `
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login TEXT UNIQUE NOT NULL,
    password BLOB NOT NULL
);

DROP TABLE IF EXISTS sessions;
CREATE TABLE sessions (
    id TEXT NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL
);

DROP TABLE IF EXISTS wallets;
CREATE TABLE wallets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER UNIQUE,
    title TEXT NOT NULL,
    general INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

DROP TABLE IF EXISTS operations;
CREATE TABLE operations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wallet_id INTEGER,
    type TEXT NOT NULL,
    amount INTEGER NOT NULL,
    date TEXT NOT NULL,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id)
);`

func ConnnectSQLite() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", "database.db")
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	db.MustExec(schema)

	tx := db.MustBegin()
	// tx.MustExec("INSERT INTO users (login, password) VALUES ($1, $2)", "Jason", "123")
	// tx.MustExec("INSERT INTO wallets (user_id, title, general) VALUES ($1, $2, $3)", "1", "my wallet", 15000)
	// tx.MustExec("INSERT INTO operations (wallet_id, type, amount, date) VALUES ($1, $2, $3, $4)", "1", "income", 1500, time.Now())
	// tx.MustExec("INSERT INTO operations (wallet_id, type, amount, date) VALUES ($1, $2, $3, $4)", "1", "expense", 1500, time.Now())
	tx.Commit()

	return db, nil
}
