package clusters

type Cluster struct {
	Id         *string           `json:"id,omitempty"`
	Name       *string           `json:"name,omitempty"`
	Properties ClusterProperties `json:"properties"`
	Sku        Sku               `json:"sku"`
	Type       *string           `json:"type,omitempty"`
}
