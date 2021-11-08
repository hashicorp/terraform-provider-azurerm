package cognitiveservicesaccounts

type VirtualNetworkRule struct {
	Id                               string  `json:"id"`
	IgnoreMissingVnetServiceEndpoint *bool   `json:"ignoreMissingVnetServiceEndpoint,omitempty"`
	State                            *string `json:"state,omitempty"`
}
