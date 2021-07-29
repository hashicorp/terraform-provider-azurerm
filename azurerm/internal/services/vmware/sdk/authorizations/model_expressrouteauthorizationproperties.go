package authorizations

type ExpressRouteAuthorizationProperties struct {
	ExpressRouteAuthorizationId  *string                                     `json:"expressRouteAuthorizationId,omitempty"`
	ExpressRouteAuthorizationKey *string                                     `json:"expressRouteAuthorizationKey,omitempty"`
	ProvisioningState            *ExpressRouteAuthorizationProvisioningState `json:"provisioningState,omitempty"`
}
