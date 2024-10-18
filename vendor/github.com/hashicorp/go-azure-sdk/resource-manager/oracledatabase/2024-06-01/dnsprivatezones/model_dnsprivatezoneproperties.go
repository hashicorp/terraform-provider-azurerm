package dnsprivatezones

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsPrivateZoneProperties struct {
	IsProtected       *bool                          `json:"isProtected,omitempty"`
	LifecycleState    *DnsPrivateZonesLifecycleState `json:"lifecycleState,omitempty"`
	Ocid              *string                        `json:"ocid,omitempty"`
	ProvisioningState *ResourceProvisioningState     `json:"provisioningState,omitempty"`
	Self              *string                        `json:"self,omitempty"`
	Serial            *int64                         `json:"serial,omitempty"`
	TimeCreated       *string                        `json:"timeCreated,omitempty"`
	Version           *string                        `json:"version,omitempty"`
	ViewId            *string                        `json:"viewId,omitempty"`
	ZoneType          *ZoneType                      `json:"zoneType,omitempty"`
}

func (o *DnsPrivateZoneProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *DnsPrivateZoneProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}
