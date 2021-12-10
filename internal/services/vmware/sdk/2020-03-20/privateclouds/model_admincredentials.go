package privateclouds

type AdminCredentials struct {
	NsxtPassword    *string `json:"nsxtPassword,omitempty"`
	NsxtUsername    *string `json:"nsxtUsername,omitempty"`
	VcenterPassword *string `json:"vcenterPassword,omitempty"`
	VcenterUsername *string `json:"vcenterUsername,omitempty"`
}
