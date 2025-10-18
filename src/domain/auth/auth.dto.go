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
