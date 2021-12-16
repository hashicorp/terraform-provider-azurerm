package databases

type DatabaseProperties struct {
	ClientProtocol    *Protocol          `json:"clientProtocol,omitempty"`
	ClusteringPolicy  *ClusteringPolicy  `json:"clusteringPolicy,omitempty"`
	EvictionPolicy    *EvictionPolicy    `json:"evictionPolicy,omitempty"`
	Modules           *[]Module          `json:"modules,omitempty"`
	Persistence       *Persistence       `json:"persistence,omitempty"`
	Port              *int64             `json:"port,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	ResourceState     *ResourceState     `json:"resourceState,omitempty"`
}
