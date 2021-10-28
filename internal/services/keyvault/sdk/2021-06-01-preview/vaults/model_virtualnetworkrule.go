package vaults

type VirtualNetworkRule struct {
	Id                               string `json:"id"`
	IgnoreMissingVnetServiceEndpoint *bool  `json:"ignoreMissingVnetServiceEndpoint,omitempty"`
}
