package privateclouds

type Circuit struct {
	ExpressRouteID               *string `json:"expressRouteID,omitempty"`
	ExpressRoutePrivatePeeringID *string `json:"expressRoutePrivatePeeringID,omitempty"`
	PrimarySubnet                *string `json:"primarySubnet,omitempty"`
	SecondarySubnet              *string `json:"secondarySubnet,omitempty"`
}
