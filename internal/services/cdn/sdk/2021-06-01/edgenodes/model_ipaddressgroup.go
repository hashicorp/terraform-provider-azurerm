package edgenodes

type IpAddressGroup struct {
	DeliveryRegion *string          `json:"deliveryRegion,omitempty"`
	Ipv4Addresses  *[]CidrIpAddress `json:"ipv4Addresses,omitempty"`
	Ipv6Addresses  *[]CidrIpAddress `json:"ipv6Addresses,omitempty"`
}
