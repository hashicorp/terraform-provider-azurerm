package labplan

type SupportInfo struct {
	Email        *string `json:"email,omitempty"`
	Instructions *string `json:"instructions,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Url          *string `json:"url,omitempty"`
}
