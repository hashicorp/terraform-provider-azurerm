package virtualnetworks

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkProperties struct {
	AllowedSubnets             *[]Subnet         `json:"allowedSubnets,omitempty"`
	CreatedDate                *string           `json:"createdDate,omitempty"`
	Description                *string           `json:"description,omitempty"`
	ExternalProviderResourceId *string           `json:"externalProviderResourceId,omitempty"`
	ExternalSubnets            *[]ExternalSubnet `json:"externalSubnets,omitempty"`
	ProvisioningState          *string           `json:"provisioningState,omitempty"`
	SubnetOverrides            *[]SubnetOverride `json:"subnetOverrides,omitempty"`
	UniqueIdentifier           *string           `json:"uniqueIdentifier,omitempty"`
}

func (o *VirtualNetworkProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *VirtualNetworkProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}
