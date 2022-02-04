package tenants

type UpdateTenantProperties struct {
	BillingConfig *BillingConfig `json:"billingConfig,omitempty"`
}
