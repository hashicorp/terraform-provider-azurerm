package diskpools

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskPoolProperties struct {
	AdditionalCapabilities *[]string          `json:"additionalCapabilities,omitempty"`
	AvailabilityZones      zones.Schema       `json:"availabilityZones"`
	Disks                  *[]Disk            `json:"disks,omitempty"`
	ProvisioningState      ProvisioningStates `json:"provisioningState"`
	Status                 OperationalStatus  `json:"status"`
	SubnetId               string             `json:"subnetId"`
}
