package signalr

type SharedPrivateLinkResourceProperties struct {
	GroupId               string                           `json:"groupId"`
	PrivateLinkResourceId string                           `json:"privateLinkResourceId"`
	ProvisioningState     *ProvisioningState               `json:"provisioningState,omitempty"`
	RequestMessage        *string                          `json:"requestMessage,omitempty"`
	Status                *SharedPrivateLinkResourceStatus `json:"status,omitempty"`
}
