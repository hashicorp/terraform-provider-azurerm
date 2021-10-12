package lab

type LabNetworkProfile struct {
	LoadBalancerId *string `json:"loadBalancerId,omitempty"`
	PublicIpId     *string `json:"publicIpId,omitempty"`
	SubnetId       *string `json:"subnetId,omitempty"`
}
