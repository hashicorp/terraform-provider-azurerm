package tenants

type BillingConfig struct {
	BillingType           *BillingType `json:"billingType,omitempty"`
	EffectiveStartDateUtc *string      `json:"effectiveStartDateUtc,omitempty"`
}
