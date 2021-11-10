package applicationtypeversion

type ApplicationTypeVersionResourceProperties struct {
	AppPackageUrl     string  `json:"appPackageUrl"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
}
