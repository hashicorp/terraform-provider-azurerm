package parse

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MaintenanceAssignmentId struct {
	TargetResourceId
	Name string
}

func MaintenanceAssignmentID(input string) (*MaintenanceAssignmentId, error) {
	// two types of ID:
	// 1: /subscriptions/<sub1>/resourcegroups/<grp1>/providers/microsoft.compute/virtualmachines/<resource1>/providers/Microsoft.Maintenance/configurationAssignments/<assignName>
	// 2: /subscriptions/<sub1>/resourcegroups/<grp1>/providers/<provider1>/<resourceParentType1>/<resourceParentName1>/<resourceType1>/<resource1>/providers/Microsoft.Maintenance/configurationAssignments/<assignName>
	groups := regexp.MustCompile(`^(.+)/providers/Microsoft\.Maintenance/configurationAssignments/([^/]+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("parsing Maintenance Assignment ID: %q: Expected 3 groups", input)
	}

	targetResourceId, name := groups[1], groups[2]
	id, err := TargetResourceID(targetResourceId)
	if err != nil {
		return nil, err
	}

	return &MaintenanceAssignmentId{
		TargetResourceId: id,
		Name:             name,
	}, nil
}

type TargetResourceId interface {
	ID() string
}

type ScopeResource struct {
	id               string
	ResourceGroup    string
	ResourceProvider string
	ResourceType     string
	ResourceName     string
}

func (r ScopeResource) ID() string {
	return r.id
}

type ScopeInResource struct {
	id                 string
	ResourceGroup      string
	ResourceProvider   string
	ResourceParentType string
	ResourceParentName string
	ResourceType       string
	ResourceName       string
}

func (r ScopeInResource) ID() string {
	return r.id
}

func TargetResourceID(input string) (TargetResourceId, error) {
	// two types of ID:
	// 1: /subscriptions/<sub1>/resourcegroups/<grp1>/providers/<provider1>/<resourceType1>/<resource1>
	// 2: /subscriptions/<sub1>/resourcegroups/<grp1>/providers/<provider1>/<resourceParentType1>/<resourceParentName1>/<resourceType1>/<resource1>

	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing target resource id %q: %+v", input, err)
	}

	if len(id.Path) != 1 && len(id.Path) != 2 {
		return nil, fmt.Errorf("parsing target resource id %q", input)
	}

	if len(id.Path) == 1 {
		var resourceType, resourceName string
		// assume there is only one left key value pair
		for k, v := range id.Path {
			resourceType = k
			resourceName = v
		}
		return ScopeResource{
			id:               input,
			ResourceGroup:    id.ResourceGroup,
			ResourceProvider: id.Provider,
			ResourceType:     resourceType,
			ResourceName:     resourceName,
		}, nil
	}

	resourceId := strings.TrimPrefix(input, "/")
	resourceId = strings.TrimSuffix(resourceId, "/")
	groups := regexp.MustCompile(`^subscriptions/.+/resource[gG]roups/.+/providers/.+/(.+)/(.+)/(.+)/(.+)$`).FindStringSubmatch(resourceId)
	if len(groups) != 5 {
		return nil, fmt.Errorf("parsing target resource id: %q", resourceId)
	}

	return ScopeInResource{
		id:                 input,
		ResourceGroup:      id.ResourceGroup,
		ResourceProvider:   id.Provider,
		ResourceParentType: groups[1],
		ResourceParentName: groups[2],
		ResourceType:       groups[3],
		ResourceName:       groups[4],
	}, nil
}
