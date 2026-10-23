package virtualnetworkaddresses

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkAddressProperties struct {
	Domain            *string                              `json:"domain,omitempty"`
	IPAddress         *string                              `json:"ipAddress,omitempty"`
	LifecycleDetails  *string                              `json:"lifecycleDetails,omitempty"`
	LifecycleState    *VirtualNetworkAddressLifecycleState `json:"lifecycleState,omitempty"`
	Ocid              *string                              `json:"ocid,omitempty"`
	ProvisioningState *AzureResourceProvisioningState      `json:"provisioningState,omitempty"`
	TimeAssigned      *string                              `json:"timeAssigned,omitempty"`
	VMOcid            *string                              `json:"vmOcid,omitempty"`
}

func (o *VirtualNetworkAddressProperties) GetTimeAssignedAsTime() (*time.Time, error) {
	if o.TimeAssigned == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeAssigned, "2006-01-02T15:04:05Z07:00")
}

func (o *VirtualNetworkAddressProperties) SetTimeAssignedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeAssigned = &formatted
}
