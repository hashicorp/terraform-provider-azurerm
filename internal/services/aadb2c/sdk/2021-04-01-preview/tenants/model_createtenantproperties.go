package tenants

type CreateTenantProperties struct {
	CountryCode string `json:"countryCode"`
	DisplayName string `json:"displayName"`
}
