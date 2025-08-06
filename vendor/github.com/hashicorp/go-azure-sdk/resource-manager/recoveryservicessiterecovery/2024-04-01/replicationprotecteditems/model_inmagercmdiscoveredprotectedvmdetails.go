package replicationprotecteditems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageRcmDiscoveredProtectedVMDetails struct {
	CreatedTimestamp       *string   `json:"createdTimestamp,omitempty"`
	DataStores             *[]string `json:"datastores,omitempty"`
	IPAddresses            *[]string `json:"ipAddresses,omitempty"`
	IsDeleted              *bool     `json:"isDeleted,omitempty"`
	LastDiscoveryTimeInUtc *string   `json:"lastDiscoveryTimeInUtc,omitempty"`
	OsName                 *string   `json:"osName,omitempty"`
	PowerStatus            *string   `json:"powerStatus,omitempty"`
	UpdatedTimestamp       *string   `json:"updatedTimestamp,omitempty"`
	VCenterFqdn            *string   `json:"vCenterFqdn,omitempty"`
	VCenterId              *string   `json:"vCenterId,omitempty"`
	VMFqdn                 *string   `json:"vmFqdn,omitempty"`
	VMwareToolsStatus      *string   `json:"vmwareToolsStatus,omitempty"`
}

func (o *InMageRcmDiscoveredProtectedVMDetails) GetCreatedTimestampAsTime() (*time.Time, error) {
	if o.CreatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *InMageRcmDiscoveredProtectedVMDetails) SetCreatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTimestamp = &formatted
}

func (o *InMageRcmDiscoveredProtectedVMDetails) GetLastDiscoveryTimeInUtcAsTime() (*time.Time, error) {
	if o.LastDiscoveryTimeInUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastDiscoveryTimeInUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *InMageRcmDiscoveredProtectedVMDetails) SetLastDiscoveryTimeInUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastDiscoveryTimeInUtc = &formatted
}

func (o *InMageRcmDiscoveredProtectedVMDetails) GetUpdatedTimestampAsTime() (*time.Time, error) {
	if o.UpdatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *InMageRcmDiscoveredProtectedVMDetails) SetUpdatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedTimestamp = &formatted
}
