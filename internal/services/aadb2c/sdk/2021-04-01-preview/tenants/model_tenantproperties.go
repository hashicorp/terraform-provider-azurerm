package tenants

type TenantProperties struct {
	BillingConfig *BillingConfig `json:"billingConfig,omitempty"`
	CountryCode   *string        `json:"countryCode,omitempty"`
	DisplayName   *string        `json:"displayName,omitempty"`
	TenantId      *string        `json:"tenantId,omitempty"`
}
