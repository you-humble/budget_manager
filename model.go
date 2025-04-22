package main

import "time"

type Wallet struct {
	ID         int64       `json:"id" db:"id"`
	UserID     int64       `json:"user_id" db:"user_id"`
	Title      string      `json:"title" db:"title"`
	General    int         `json:"general" db:"general"`
	Operations []Operation `json:"operations"`
}

type Operation struct {
	ID       int64     `json:"id" db:"id"`
	WalletID int64     `json:"wallet_id" db:"wallet_id"`
	Type     string    `json:"type" db:"type"`
	Amount   int       `json:"amount" db:"amount"`
	Date     time.Time `json:"date" db:"date"`
}
