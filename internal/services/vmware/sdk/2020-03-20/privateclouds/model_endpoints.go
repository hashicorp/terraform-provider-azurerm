package privateclouds

type Endpoints struct {
	HcxCloudManager *string `json:"hcxCloudManager,omitempty"`
	NsxtManager     *string `json:"nsxtManager,omitempty"`
	Vcsa            *string `json:"vcsa,omitempty"`
}
