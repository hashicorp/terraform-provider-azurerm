package cognitiveservicesaccounts

type QuotaLimit struct {
	Count         *float64          `json:"count,omitempty"`
	RenewalPeriod *float64          `json:"renewalPeriod,omitempty"`
	Rules         *[]ThrottlingRule `json:"rules,omitempty"`
}
