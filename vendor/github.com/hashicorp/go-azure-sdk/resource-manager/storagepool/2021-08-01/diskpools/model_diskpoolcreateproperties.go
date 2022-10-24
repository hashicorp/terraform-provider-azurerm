package diskpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskPoolCreateProperties struct {
	AdditionalCapabilities *[]string `json:"additionalCapabilities,omitempty"`
	AvailabilityZones      *[]string `json:"availabilityZones,omitempty"`
	Disks                  *[]Disk   `json:"disks,omitempty"`
	SubnetId               string    `json:"subnetId"`
}
