package dnsprivateviews

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsPrivateViewProperties struct {
	DisplayName       *string                        `json:"displayName,omitempty"`
	IsProtected       *bool                          `json:"isProtected,omitempty"`
	LifecycleState    *DnsPrivateViewsLifecycleState `json:"lifecycleState,omitempty"`
	Ocid              *string                        `json:"ocid,omitempty"`
	ProvisioningState *ResourceProvisioningState     `json:"provisioningState,omitempty"`
	Self              *string                        `json:"self,omitempty"`
	TimeCreated       *string                        `json:"timeCreated,omitempty"`
	TimeUpdated       *string                        `json:"timeUpdated,omitempty"`
}

func (o *DnsPrivateViewProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *DnsPrivateViewProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}

func (o *DnsPrivateViewProperties) GetTimeUpdatedAsTime() (*time.Time, error) {
	if o.TimeUpdated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeUpdated, "2006-01-02T15:04:05Z07:00")
}

func (o *DnsPrivateViewProperties) SetTimeUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeUpdated = &formatted
}
