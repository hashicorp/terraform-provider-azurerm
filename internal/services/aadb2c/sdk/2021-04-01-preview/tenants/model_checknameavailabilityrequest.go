package tenants

type CheckNameAvailabilityRequest struct {
	CountryCode *string `json:"countryCode,omitempty"`
	Name        *string `json:"name,omitempty"`
}
