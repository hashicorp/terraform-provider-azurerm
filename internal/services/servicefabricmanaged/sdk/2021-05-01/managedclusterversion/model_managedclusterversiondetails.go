package managedclusterversion

type ManagedClusterVersionDetails struct {
	ClusterCodeVersion *string `json:"clusterCodeVersion,omitempty"`
	OsType             *OsType `json:"osType,omitempty"`
	SupportExpiryUtc   *string `json:"supportExpiryUtc,omitempty"`
}
