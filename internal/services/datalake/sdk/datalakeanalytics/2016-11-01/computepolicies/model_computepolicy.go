package computepolicies

type ComputePolicy struct {
	Id         *string                  `json:"id,omitempty"`
	Name       *string                  `json:"name,omitempty"`
	Properties *ComputePolicyProperties `json:"properties,omitempty"`
	Type       *string                  `json:"type,omitempty"`
}
