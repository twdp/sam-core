package req

type EmailLoginDto struct {
	Email string `json:"email"`

	Password string `json:"password"`
}
