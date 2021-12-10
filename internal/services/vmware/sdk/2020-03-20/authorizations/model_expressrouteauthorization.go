package authorizations

type ExpressRouteAuthorization struct {
	Id         *string                              `json:"id,omitempty"`
	Name       *string                              `json:"name,omitempty"`
	Properties *ExpressRouteAuthorizationProperties `json:"properties,omitempty"`
	Type       *string                              `json:"type,omitempty"`
}
