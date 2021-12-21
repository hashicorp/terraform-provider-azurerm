package cognitiveservicesaccounts

type ThrottlingRule struct {
	Count                    *float64               `json:"count,omitempty"`
	DynamicThrottlingEnabled *bool                  `json:"dynamicThrottlingEnabled,omitempty"`
	Key                      *string                `json:"key,omitempty"`
	MatchPatterns            *[]RequestMatchPattern `json:"matchPatterns,omitempty"`
	MinCount                 *float64               `json:"minCount,omitempty"`
	RenewalPeriod            *float64               `json:"renewalPeriod,omitempty"`
}
