package resourceskus

type ResourceSkuInfo struct {
	ApiVersion   *string                    `json:"apiVersion,omitempty"`
	Capabilities *[]ResourceSkuCapability   `json:"capabilities,omitempty"`
	LocationInfo *ResourceSkuLocationInfo   `json:"locationInfo,omitempty"`
	Name         *string                    `json:"name,omitempty"`
	ResourceType *string                    `json:"resourceType,omitempty"`
	Restrictions *[]ResourceSkuRestrictions `json:"restrictions,omitempty"`
	Tier         *string                    `json:"tier,omitempty"`
}
