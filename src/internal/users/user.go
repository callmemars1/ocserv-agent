package users

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsBanned bool   `json:"is_banned"`
}
