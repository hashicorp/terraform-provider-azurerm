package redisenterprise

type ClusterUpdate struct {
	Properties *ClusterProperties `json:"properties,omitempty"`
	Sku        *Sku               `json:"sku,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
