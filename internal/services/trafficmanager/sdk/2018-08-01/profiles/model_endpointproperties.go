package profiles

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)



type EndpointProperties struct {
	CustomHeaders *[] `json:"customHeaders,omitempty"`
	EndpointLocation *string `json:"endpointLocation,omitempty"`
	EndpointMonitorStatus *EndpointMonitorStatus `json:"endpointMonitorStatus,omitempty"`
	EndpointStatus *EndpointStatus `json:"endpointStatus,omitempty"`
	GeoMapping *[]string `json:"geoMapping,omitempty"`
	MinChildEndpoints *int64 `json:"minChildEndpoints,omitempty"`
	MinChildEndpointsIPv4 *int64 `json:"minChildEndpointsIPv4,omitempty"`
	MinChildEndpointsIPv6 *int64 `json:"minChildEndpointsIPv6,omitempty"`
	Priority *int64 `json:"priority,omitempty"`
	Subnets *[] `json:"subnets,omitempty"`
	Target *string `json:"target,omitempty"`
	TargetResourceId *string `json:"targetResourceId,omitempty"`
	Weight *int64 `json:"weight,omitempty"`
}







