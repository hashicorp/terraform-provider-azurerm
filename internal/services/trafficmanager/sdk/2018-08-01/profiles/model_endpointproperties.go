package profiles

type EndpointProperties struct {
	CustomHeaders         *[]EndpointPropertiesCustomHeadersInlined `json:"customHeaders,omitempty"`
	EndpointLocation      *string                                   `json:"endpointLocation,omitempty"`
	EndpointMonitorStatus *EndpointMonitorStatus                    `json:"endpointMonitorStatus,omitempty"`
	EndpointStatus        *EndpointStatus                           `json:"endpointStatus,omitempty"`
	GeoMapping            *[]string                                 `json:"geoMapping,omitempty"`
	MinChildEndpoints     *int64                                    `json:"minChildEndpoints,omitempty"`
	MinChildEndpointsIPv4 *int64                                    `json:"minChildEndpointsIPv4,omitempty"`
	MinChildEndpointsIPv6 *int64                                    `json:"minChildEndpointsIPv6,omitempty"`
	Priority              *int64                                    `json:"priority,omitempty"`
	Subnets               *[]EndpointPropertiesSubnetsInlined       `json:"subnets,omitempty"`
	Target                *string                                   `json:"target,omitempty"`
	TargetResourceId      *string                                   `json:"targetResourceId,omitempty"`
	Weight                *int64                                    `json:"weight,omitempty"`
}
