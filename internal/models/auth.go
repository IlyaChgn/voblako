package models

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type FullUserData struct {
	User
	SessionID string
}

type AuthData struct {
	User   User `json:"user"`
	IsAuth bool `json:"is_auth"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupData struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}
