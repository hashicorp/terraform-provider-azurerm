package clusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterNode struct {
	CoreCount                 *float64                   `json:"coreCount,omitempty"`
	EhcResourceId             *string                    `json:"ehcResourceId,omitempty"`
	Id                        *float64                   `json:"id,omitempty"`
	LastLicensingTimestamp    *string                    `json:"lastLicensingTimestamp,omitempty"`
	Manufacturer              *string                    `json:"manufacturer,omitempty"`
	MemoryInGiB               *float64                   `json:"memoryInGiB,omitempty"`
	Model                     *string                    `json:"model,omitempty"`
	Name                      *string                    `json:"name,omitempty"`
	NodeType                  *ClusterNodeType           `json:"nodeType,omitempty"`
	OemActivation             *OemActivation             `json:"oemActivation,omitempty"`
	OsDisplayVersion          *string                    `json:"osDisplayVersion,omitempty"`
	OsName                    *string                    `json:"osName,omitempty"`
	OsVersion                 *string                    `json:"osVersion,omitempty"`
	SerialNumber              *string                    `json:"serialNumber,omitempty"`
	WindowsServerSubscription *WindowsServerSubscription `json:"windowsServerSubscription,omitempty"`
}

func (o *ClusterNode) GetLastLicensingTimestampAsTime() (*time.Time, error) {
	if o.LastLicensingTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastLicensingTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ClusterNode) SetLastLicensingTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastLicensingTimestamp = &formatted
}
