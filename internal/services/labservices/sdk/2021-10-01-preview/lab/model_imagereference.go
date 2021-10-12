package lab

type ImageReference struct {
	ExactVersion *string `json:"exactVersion,omitempty"`
	Id           *string `json:"id,omitempty"`
	Offer        *string `json:"offer,omitempty"`
	Publisher    *string `json:"publisher,omitempty"`
	Sku          *string `json:"sku,omitempty"`
	Version      *string `json:"version,omitempty"`
}
