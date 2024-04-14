package models

type RefUser struct {
	RefUserId uint   `json:"user_id"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Token     string `json:"token"`
	Is_Admin  bool   `json:"is_admin"`
}
