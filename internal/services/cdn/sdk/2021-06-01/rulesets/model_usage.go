package rulesets

type Usage struct {
	CurrentValue int64     `json:"currentValue"`
	Id           *string   `json:"id,omitempty"`
	Limit        int64     `json:"limit"`
	Name         UsageName `json:"name"`
	Unit         UsageUnit `json:"unit"`
}
