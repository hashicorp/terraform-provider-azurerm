package privateclouds

type PrivateCloudUpdateProperties struct {
	IdentitySources   *[]IdentitySource  `json:"identitySources,omitempty"`
	Internet          *InternetEnum      `json:"internet,omitempty"`
	ManagementCluster *ManagementCluster `json:"managementCluster,omitempty"`
}
