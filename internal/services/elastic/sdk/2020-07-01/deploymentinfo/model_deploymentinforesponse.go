package deploymentinfo

type DeploymentInfoResponse struct {
	DiskCapacity   *string                  `json:"diskCapacity,omitempty"`
	MemoryCapacity *string                  `json:"memoryCapacity,omitempty"`
	Status         *ElasticDeploymentStatus `json:"status,omitempty"`
	Version        *string                  `json:"version,omitempty"`
}
