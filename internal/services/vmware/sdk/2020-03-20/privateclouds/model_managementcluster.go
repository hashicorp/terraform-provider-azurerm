package privateclouds

type ManagementCluster struct {
	ClusterId         *int64                    `json:"clusterId,omitempty"`
	ClusterSize       int64                     `json:"clusterSize"`
	Hosts             *[]string                 `json:"hosts,omitempty"`
	ProvisioningState *ClusterProvisioningState `json:"provisioningState,omitempty"`
}
