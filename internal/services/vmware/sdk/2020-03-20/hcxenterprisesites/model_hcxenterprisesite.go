package hcxenterprisesites

type HcxEnterpriseSite struct {
	Id         *string                      `json:"id,omitempty"`
	Name       *string                      `json:"name,omitempty"`
	Properties *HcxEnterpriseSiteProperties `json:"properties,omitempty"`
	Type       *string                      `json:"type,omitempty"`
}
