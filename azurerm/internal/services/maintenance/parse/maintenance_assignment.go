package parse

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MaintenanceAssignmentId struct {
	*TargetResourceId
	ResourceId string
	Name       string
}

func MaintenanceAssignmentID(input string) (*MaintenanceAssignmentId, error) {
	// two types of ID:
	// 1: /subscriptions/<sub1>/resourcegroups/<grp1>/providers/microsoft.compute/virtualmachines/<resource1>/providers/Microsoft.Maintenance/configurationAssignments/<assignName>
	// 2: /subscriptions/<sub1>/resourcegroups/<grp1>/providers/<provider1>/<resourceParentType1>/<resourceParentName1>/<resourceType1>/<resource1>/providers/Microsoft.Maintenance/configurationAssignments/<assignName>
	groups := regexp.MustCompile(`^(.+)/providers/Microsoft\.Maintenance/configurationAssignments/([^/]+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("parsing Maintenance Assignment ID: %q", input)
	}

	targetResourceId, name := groups[1], groups[2]
	id, err := TargetResourceID(targetResourceId)
	if err != nil {
		return nil, err
	}

	return &MaintenanceAssignmentId{
		TargetResourceId: id,
		ResourceId:       targetResourceId,
		Name:             name,
	}, nil
}

type TargetResourceId struct {
	HasParentResource  bool
	ResourceGroup      string
	ResourceProvider   string
	ResourceParentType string
	ResourceParentName string
	ResourceType       string
	ResourceName       string
}

func TargetResourceID(input string) (*TargetResourceId, error) {
	// two types of ID:
	// 1: /subscriptions/<sub1>/resourcegroups/<grp1>/providers/<provider1>/<resourceType1>/<resource1>
	// 2: /subscriptions/<sub1>/resourcegroups/<grp1>/providers/<provider1>/<resourceParentType1>/<resourceParentName1>/<resourceType1>/<resource1>

	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing target resource id %q: %+v", input, err)
	}

	var resourceParentType, resourceParentName, resourceType, resourceName string
	var hasParentResource bool
	if len(id.Path) != 1 && len(id.Path) != 2 {
		return nil, fmt.Errorf("parsing target resource id %q", input)
	}

	if len(id.Path) == 1 {
		// assume there is only one left key value pair
		for k, v := range id.Path {
			resourceType = k
			resourceName = v
		}
	} else {
		hasParentResource = true

		input = strings.TrimPrefix(input, "/")
		input = strings.TrimSuffix(input, "/")
		groups := regexp.MustCompile(`^subscriptions/.+/resource[gG]roups/.+/providers/.+/(.+)/(.+)/(.+)/(.+)$`).FindStringSubmatch(input)
		if len(groups) != 5 {
			return nil, fmt.Errorf("parsing target resource id: %q", input)
		}

		resourceParentType = groups[1]
		resourceParentName = groups[2]
		resourceType = groups[3]
		resourceName = groups[4]
	}

	return &TargetResourceId{
		ResourceGroup:      id.ResourceGroup,
		ResourceProvider:   id.Provider,
		ResourceParentType: resourceParentType,
		ResourceParentName: resourceParentName,
		ResourceType:       resourceType,
		ResourceName:       resourceName,
		HasParentResource:  hasParentResource,
	}, nil
}
