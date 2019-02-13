package req

type EmailLoginDto struct {
	Email string `json:"email"`

	Password string `json:"password"`

	// 哪个端
	Terminal int8
}
