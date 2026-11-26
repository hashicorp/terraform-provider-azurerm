package expressroutecrossconnections

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateExpressRouteCircuitHealthCheckResponse struct {
	Details               *MigrateExpressRouteCircuitHealthCheckDetails `json:"details,omitempty"`
	FailureReason         *string                                       `json:"failureReason,omitempty"`
	NewCrossConnectionURL *string                                       `json:"newCrossConnectionUrl,omitempty"`
	NewSTag               *string                                       `json:"newSTag,omitempty"`
	Phase                 *string                                       `json:"phase,omitempty"`
	PrepareExpiryTime     *string                                       `json:"prepareExpiryTime,omitempty"`
	PreparedAt            *string                                       `json:"preparedAt,omitempty"`
	ShouldRollback        *bool                                         `json:"shouldRollback,omitempty"`
	Status                *string                                       `json:"status,omitempty"`
}

func (o *MigrateExpressRouteCircuitHealthCheckResponse) GetPrepareExpiryTimeAsTime() (*time.Time, error) {
	if o.PrepareExpiryTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PrepareExpiryTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MigrateExpressRouteCircuitHealthCheckResponse) SetPrepareExpiryTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PrepareExpiryTime = &formatted
}

func (o *MigrateExpressRouteCircuitHealthCheckResponse) GetPreparedAtAsTime() (*time.Time, error) {
	if o.PreparedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PreparedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *MigrateExpressRouteCircuitHealthCheckResponse) SetPreparedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PreparedAt = &formatted
}
