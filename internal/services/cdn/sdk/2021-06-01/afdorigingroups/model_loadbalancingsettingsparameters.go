package afdorigingroups

type LoadBalancingSettingsParameters struct {
	AdditionalLatencyInMilliseconds *int64 `json:"additionalLatencyInMilliseconds,omitempty"`
	SampleSize                      *int64 `json:"sampleSize,omitempty"`
	SuccessfulSamplesRequired       *int64 `json:"successfulSamplesRequired,omitempty"`
}
