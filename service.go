package main

import (
	"log/slog"
	"strings"
	"time"
)

type WalletRepository interface {
	CreateWallet(w Wallet) (int64, error)
	AddOperation(userID int64, op Operation) error
	WalletByID(id int64) (Wallet, error)
}

type walletService struct {
	repo WalletRepository
}

func NewWalletService(repo WalletRepository) *walletService {
	return &walletService{repo: repo}
}

func (s *walletService) CreateWallet(w Wallet) (Wallet, error) {
	id, err := s.repo.CreateWallet(w)
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
