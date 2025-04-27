package user

type User struct {
	ID       int64  `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password []byte `json:"password" db:"password"`
}
