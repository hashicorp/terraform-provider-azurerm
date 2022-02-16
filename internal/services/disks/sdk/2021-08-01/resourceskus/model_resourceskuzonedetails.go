package resourceskus

type ResourceSkuZoneDetails struct {
	Capabilities *[]ResourceSkuCapability `json:"capabilities,omitempty"`
	Name         *[]string                `json:"name,omitempty"`
}
