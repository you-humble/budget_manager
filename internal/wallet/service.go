package wallet

import (
	"log/slog"
	"strings"
	"time"
)

type Repository interface {
	Save(w Wallet) (int64, error)
	AddOperation(userID int64, op Operation) error
	WalletByID(id int64) (Wallet, error)
}

type walletService struct {
	repo Repository
}

func NewService(repo Repository) *walletService {
	return &walletService{repo: repo}
}

func (s *walletService) Save(w Wallet) (Wallet, error) {
	id, err := s.repo.Save(w)
	if err != nil {
		slog.Error("failed to create wallet", slog.String("error", err.Error()))
		return Wallet{}, err
	}

	w.ID = id
	w.Operations = make([]Operation, 0, 10)
	return w, nil
}

// TODo: prohibit a negative general
func (s *walletService) AddOperation(userID int64, op Operation) error {
	op.Type = strings.ToLower(strings.TrimSpace(op.Type))
	op.Date = time.Now().Local().UTC()

	if err := s.repo.AddOperation(userID, op); err != nil {
		slog.Error("failed to create add operation", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *walletService) ShowWallet(userID int64) (Wallet, error) {
	return s.repo.WalletByID(userID)
}
