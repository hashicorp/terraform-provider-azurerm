package accounts

type UpdateComputePolicyWithAccountParameters struct {
	Name       string                         `json:"name"`
	Properties *UpdateComputePolicyProperties `json:"properties,omitempty"`
}
