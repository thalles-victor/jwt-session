package domain_auth

type SignUpDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequestRecoveryDto struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
	Code        string `json:"code"`
}
