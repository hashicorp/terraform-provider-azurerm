package parse

import (
	"fmt"
	"regexp"

	parseCompute "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
)

type MaintenanceAssignmentDedicatedHostId struct {
	DedicatedHostId    *parseCompute.DedicatedHostId
	DedicatedHostIdRaw string
	Name               string
}

func MaintenanceAssignmentDedicatedHostID(input string) (*MaintenanceAssignmentDedicatedHostId, error) {
	groups := regexp.MustCompile(`^(.+)/providers/Microsoft\.Maintenance/configurationAssignments/([^/]+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("parsing Maintenance Assignment Dedicated Host ID (%q)", input)
	}

	targetResourceId, name := groups[1], groups[2]
	dedicatedHostID, err := parseCompute.DedicatedHostID(targetResourceId)
	if err != nil {
		return nil, fmt.Errorf("parsing Maintenance Assignment Dedicated Host ID: %q: Expected valid Dedicated Host ID", input)
	}

	return &MaintenanceAssignmentDedicatedHostId{
		DedicatedHostId:    dedicatedHostID,
		DedicatedHostIdRaw: targetResourceId,
		Name:               name,
	}, nil
}
