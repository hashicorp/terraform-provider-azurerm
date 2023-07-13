// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/dedicatedhosts"
)

type MaintenanceAssignmentDedicatedHostId struct {
	DedicatedHostId    dedicatedhosts.HostId
	DedicatedHostIdRaw string
	Name               string
}

func MaintenanceAssignmentDedicatedHostID(input string) (*MaintenanceAssignmentDedicatedHostId, error) {
	groups := regexp.MustCompile(`^(.+)/providers/Microsoft\.Maintenance/configurationAssignments/([^/]+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("parsing Maintenance Assignment Dedicated Host ID (%q)", input)
	}

	targetResourceId, name := groups[1], groups[2]
	dedicatedHostID, err := dedicatedhosts.ParseHostIDInsensitively(targetResourceId)
	if err != nil {
		return nil, fmt.Errorf("parsing Maintenance Assignment Dedicated Host ID: %q: Expected valid Dedicated Host ID", input)
	}

	return &MaintenanceAssignmentDedicatedHostId{
		DedicatedHostId:    *dedicatedHostID,
		DedicatedHostIdRaw: targetResourceId,
		Name:               name,
	}, nil
}
