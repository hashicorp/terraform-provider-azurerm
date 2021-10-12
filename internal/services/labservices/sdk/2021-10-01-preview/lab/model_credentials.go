package lab

type Credentials struct {
	Password *string `json:"password,omitempty"`
	Username string  `json:"username"`
}
