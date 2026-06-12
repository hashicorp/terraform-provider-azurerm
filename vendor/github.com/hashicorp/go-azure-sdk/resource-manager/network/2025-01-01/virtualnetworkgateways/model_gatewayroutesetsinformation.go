package virtualnetworkgateways

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayRouteSetsInformation struct {
	CircuitsMetadataMap     *map[string]CircuitMetadataMap `json:"circuitsMetadataMap,omitempty"`
	LastComputedTime        *string                        `json:"lastComputedTime,omitempty"`
	NextEligibleComputeTime *string                        `json:"nextEligibleComputeTime,omitempty"`
	RouteSetVersion         *string                        `json:"routeSetVersion,omitempty"`
	RouteSets               *[]GatewayRouteSet             `json:"routeSets,omitempty"`
}

func (o *GatewayRouteSetsInformation) GetLastComputedTimeAsTime() (*time.Time, error) {
	if o.LastComputedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastComputedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *GatewayRouteSetsInformation) SetLastComputedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastComputedTime = &formatted
}

func (o *GatewayRouteSetsInformation) GetNextEligibleComputeTimeAsTime() (*time.Time, error) {
	if o.NextEligibleComputeTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NextEligibleComputeTime, "2006-01-02T15:04:05Z07:00")
}

func (o *GatewayRouteSetsInformation) SetNextEligibleComputeTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NextEligibleComputeTime = &formatted
}
