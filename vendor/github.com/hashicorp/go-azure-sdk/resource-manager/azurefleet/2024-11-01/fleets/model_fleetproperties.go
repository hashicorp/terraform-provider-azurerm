package fleets

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetProperties struct {
	AdditionalLocationsProfile *AdditionalLocationsProfile `json:"additionalLocationsProfile,omitempty"`
	ComputeProfile             ComputeProfile              `json:"computeProfile"`
	ProvisioningState          *ProvisioningState          `json:"provisioningState,omitempty"`
	RegularPriorityProfile     *RegularPriorityProfile     `json:"regularPriorityProfile,omitempty"`
	SpotPriorityProfile        *SpotPriorityProfile        `json:"spotPriorityProfile,omitempty"`
	TimeCreated                *string                     `json:"timeCreated,omitempty"`
	UniqueId                   *string                     `json:"uniqueId,omitempty"`
	VMAttributes               *VMAttributes               `json:"vmAttributes,omitempty"`
	VMSizesProfile             []VMSizeProfile             `json:"vmSizesProfile"`
}

func (o *FleetProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *FleetProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}
