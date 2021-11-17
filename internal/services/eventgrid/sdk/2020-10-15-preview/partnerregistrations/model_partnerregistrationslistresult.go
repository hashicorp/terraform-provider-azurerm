package partnerregistrations

type PartnerRegistrationsListResult struct {
	NextLink *string                `json:"nextLink,omitempty"`
	Value    *[]PartnerRegistration `json:"value,omitempty"`
}
