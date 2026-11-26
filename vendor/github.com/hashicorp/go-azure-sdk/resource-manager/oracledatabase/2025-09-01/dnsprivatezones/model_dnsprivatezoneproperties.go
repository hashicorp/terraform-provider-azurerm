package dnsprivatezones

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsPrivateZoneProperties struct {
	IsProtected       bool                          `json:"isProtected"`
	LifecycleState    DnsPrivateZonesLifecycleState `json:"lifecycleState"`
	Ocid              string                        `json:"ocid"`
	ProvisioningState *ResourceProvisioningState    `json:"provisioningState,omitempty"`
	Self              string                        `json:"self"`
	Serial            int64                         `json:"serial"`
	TimeCreated       string                        `json:"timeCreated"`
	Version           string                        `json:"version"`
	ViewId            *string                       `json:"viewId,omitempty"`
	ZoneType          ZoneType                      `json:"zoneType"`
}

func (o *DnsPrivateZoneProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *DnsPrivateZoneProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = formatted
}
