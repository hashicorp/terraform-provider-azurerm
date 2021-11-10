package skus

type ResourceSku struct {
	Kind         *string                    `json:"kind,omitempty"`
	Locations    *[]string                  `json:"locations,omitempty"`
	Name         *string                    `json:"name,omitempty"`
	ResourceType *string                    `json:"resourceType,omitempty"`
	Restrictions *[]ResourceSkuRestrictions `json:"restrictions,omitempty"`
	Tier         *string                    `json:"tier,omitempty"`
}
