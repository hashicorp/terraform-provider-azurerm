package hdinsights

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPoolComputeProfile struct {
	AvailabilityZones *zones.Schema `json:"availabilityZones,omitempty"`
	Count             *int64        `json:"count,omitempty"`
	VMSize            string        `json:"vmSize"`
}
