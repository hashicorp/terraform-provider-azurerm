package virtualmachine

type ResetPasswordBody struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
