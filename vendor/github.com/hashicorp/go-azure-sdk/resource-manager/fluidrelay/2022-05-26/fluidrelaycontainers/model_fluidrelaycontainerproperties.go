package fluidrelaycontainers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluidRelayContainerProperties struct {
	CreationTime      *string            `json:"creationTime,omitempty"`
	FrsContainerId    *string            `json:"frsContainerId,omitempty"`
	FrsTenantId       *string            `json:"frsTenantId,omitempty"`
	LastAccessTime    *string            `json:"lastAccessTime,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}

func (o *FluidRelayContainerProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *FluidRelayContainerProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *FluidRelayContainerProperties) GetLastAccessTimeAsTime() (*time.Time, error) {
	if o.LastAccessTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastAccessTime, "2006-01-02T15:04:05Z07:00")
}

func (o *FluidRelayContainerProperties) SetLastAccessTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastAccessTime = &formatted
}
