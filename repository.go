package main

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type walletOptRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) *walletOptRepository {
	return &walletOptRepository{db: db}
}

func (repo *walletOptRepository) CreateWallet(w Wallet) (int64, error) {
	var id int64
	if err := repo.db.Get(&id, `
	INSERT INTO wallets
	(user_id, title, general)
	VALUES
	($1, $2, $3)
	RETURNING id;`, w.UserID, w.Title, w.General); err != nil {
		return 0, fmt.Errorf("CreateWallet - error after processing the INSERT query 1: %w", err)
	}

	return id, nil
}

func (repo *walletOptRepository) AddOperation(userID int64, op Operation) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return fmt.Errorf("AddOperation - failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	var general int
	if err := repo.db.Get(&general, `
	SELECT general FROM wallets WHERE user_id = $1;`,
		userID); err != nil {
		tx.Rollback()
		return fmt.Errorf("AddOperation - error after processing the SELECT query 2: %w", err)
	}

	if general-op.Amount < 0 {
		tx.Rollback()
		return fmt.Errorf("AddOperation - amount more than general")
	}

	if _, err := repo.db.Exec(`
	INSERT INTO operations
	(wallet_id, type, amount, date)
	VALUES
	($1, $2, $3, $4);`, op.WalletID, op.Type, op.Amount, op.Date); err != nil {
		tx.Rollback()
		return fmt.Errorf("AddOperation - error after processing the INSERT query 3: %w", err)
	}

	if op.Type == "income" {
		if _, err := repo.db.Exec(`
		UPDATE wallets
			SET general = general + $1
		WHERE user_id = $2;`, op.Amount, userID); err != nil {
			tx.Rollback()
			return fmt.Errorf("AddOperation - error after processing the INSERT query 4: %w", err)
		}
	} else if op.Type == "expense" {
		if _, err := repo.db.Exec(`
		UPDATE wallets
			SET general = general - $1
		WHERE user_id = $2;`, op.Amount, userID); err != nil {
			tx.Rollback()
			return fmt.Errorf("AddOperation - error after processing the INSERT query 5: %w", err)
		}
	} else {
		tx.Rollback()
		return fmt.Errorf("AddOperation - error after processing the INSERT query 6: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("AddOperation - failed to commit: %w", err)
	}
	return nil
}

func (repo *walletOptRepository) WalletByID(id int64) (Wallet, error) {
	var wwo []WalletWithOperations
	if err := repo.db.Select(&wwo,
		`SELECT
			wallets.id AS wallet_id,
			wallets.user_id AS wallet_user_id,
			wallets.title AS wallet_title,
			wallets.general AS wallet_general,
			operations.id AS operations_id,
			operations.type AS operations_type,
			operations.amount AS operations_amount,
			operations.date AS operations_date
		FROM wallets LEFT JOIN operations
			ON wallets.id = operations.wallet_id
		WHERE wallets.user_id = $1
		ORDER BY operations_date;`,
		id); err != nil {
		return Wallet{}, fmt.Errorf(
			"failed to show the wallet by id = %d due error: %w", id, err,
		)
	}

	return newWallet(wwo)
}

func newWallet(rawWallet []WalletWithOperations) (Wallet, error) {
	w := Wallet{
		ID:         rawWallet[0].WalletID,
		UserID:     rawWallet[0].WalletUserID,
		Title:      rawWallet[0].WalletTitle,
		General:    rawWallet[0].WalletGeneral,
		Operations: make([]Operation, 0, len(rawWallet)),
	}

	if rawWallet[0].OpID != nil {
		for _, wwo := range rawWallet {
			date, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", *wwo.OpDate)
			if err != nil {
				return Wallet{}, fmt.Errorf(
					"failed to parse date %q due error: %w", *wwo.OpDate, err,
				)
			}
			w.Operations = append(w.Operations, Operation{
				ID:     *wwo.OpID,
				Type:   *wwo.OpType,
				Amount: *wwo.OpAmount,
				Date:   date,
			})
		}
	}

	return w, nil
}
