package skus

type ResourceSkuRestrictionInfo struct {
	Locations *[]string `json:"locations,omitempty"`
	Zones     *[]string `json:"zones,omitempty"`
}
