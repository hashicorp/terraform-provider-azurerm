package hcxenterprisesites

type HcxEnterpriseSiteProperties struct {
	ActivationKey *string                  `json:"activationKey,omitempty"`
	Status        *HcxEnterpriseSiteStatus `json:"status,omitempty"`
}
