package servers

type GatewayDetails struct {
	DmtsClusterUri    *string `json:"dmtsClusterUri,omitempty"`
	GatewayObjectId   *string `json:"gatewayObjectId,omitempty"`
	GatewayResourceId *string `json:"gatewayResourceId,omitempty"`
}
