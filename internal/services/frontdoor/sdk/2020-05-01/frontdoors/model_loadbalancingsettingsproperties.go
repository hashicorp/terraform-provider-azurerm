package frontdoors

type LoadBalancingSettingsProperties struct {
	AdditionalLatencyMilliseconds *int64                  `json:"additionalLatencyMilliseconds,omitempty"`
	ResourceState                 *FrontDoorResourceState `json:"resourceState,omitempty"`
	SampleSize                    *int64                  `json:"sampleSize,omitempty"`
	SuccessfulSamplesRequired     *int64                  `json:"successfulSamplesRequired,omitempty"`
}
