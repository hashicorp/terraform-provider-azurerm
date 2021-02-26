package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ManagementGroupTemplateDeploymentId struct {
	ManagementGroupName string
	DeploymentName      string
}

func NewManagementGroupTemplateDeploymentID(managementGroupName, deploymentName string) ManagementGroupTemplateDeploymentId {
	return ManagementGroupTemplateDeploymentId{
		ManagementGroupName: managementGroupName,
		DeploymentName:      deploymentName,
	}
}

func (id ManagementGroupTemplateDeploymentId) String() string {
	segments := []string{
		fmt.Sprintf("Deployment Name %q", id.DeploymentName),
		fmt.Sprintf("Management Group Name %q", id.ManagementGroupName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Management Group Template Deployment", segmentsStr)
}

func (id ManagementGroupTemplateDeploymentId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Resources/deployments/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupName, id.DeploymentName)
}

// ManagementGroupTemplateDeploymentID parses a ManagementGroupTemplateDeployment ID into an ManagementGroupTemplateDeploymentId struct
func ManagementGroupTemplateDeploymentID(input string) (*ManagementGroupTemplateDeploymentId, error) {
	idURL, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components)%2 != 0 {
		return nil, fmt.Errorf("The number of path segments is not divisible by 2 in %q", path)
	}

	componentMap := make(map[string]string, len(components)/2)
	for current := 0; current < len(components); current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return nil, fmt.Errorf("Key/Value cannot be empty strings. Key: '%s', Value: '%s'", key, value)
		}
		componentMap[key] = value
	}

	// Build up a TargetResourceID from the map
	id := &azure.ResourceID{}
	id.Path = componentMap

	if provider, ok := componentMap["providers"]; ok {
		id.Provider = provider
		delete(componentMap, "providers")
	}

	resourceId := ManagementGroupTemplateDeploymentId{}

	if resourceId.ManagementGroupName, err = id.PopSegment("managementGroups"); err != nil {
		return nil, err
	}
	if resourceId.DeploymentName, err = id.PopSegment("deployments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
