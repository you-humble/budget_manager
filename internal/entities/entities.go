package entities

type WalletWithOperations struct {
	WalletID      int64   `db:"wallet_id"`
	WalletUserID  int64   `db:"wallet_user_id"`
	WalletTitle   string  `db:"wallet_title"`
	WalletGeneral int     `db:"wallet_general"`
	OpID          *int64  `db:"operations_id"`
	OpType        *string `db:"operations_type"`
	OpAmount      *int    `db:"operations_amount"`
	OpDate        *string `db:"operations_date"`
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
