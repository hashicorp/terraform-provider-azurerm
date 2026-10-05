package dnsprivateviews

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsPrivateViewProperties struct {
	DisplayName       string                        `json:"displayName"`
	IsProtected       bool                          `json:"isProtected"`
	LifecycleState    DnsPrivateViewsLifecycleState `json:"lifecycleState"`
	Ocid              string                        `json:"ocid"`
	ProvisioningState *ResourceProvisioningState    `json:"provisioningState,omitempty"`
	Self              string                        `json:"self"`
	TimeCreated       string                        `json:"timeCreated"`
	TimeUpdated       string                        `json:"timeUpdated"`
}

func (o *DnsPrivateViewProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *DnsPrivateViewProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = formatted
}

func (o *DnsPrivateViewProperties) GetTimeUpdatedAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.TimeUpdated, "2006-01-02T15:04:05Z07:00")
}

func (o *DnsPrivateViewProperties) SetTimeUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeUpdated = formatted
}
