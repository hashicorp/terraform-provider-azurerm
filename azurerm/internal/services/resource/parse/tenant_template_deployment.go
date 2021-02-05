package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TenantTemplateDeploymentId struct {
	DeploymentName string
}

func NewTenantTemplateDeploymentID(deploymentName string) TenantTemplateDeploymentId {
	return TenantTemplateDeploymentId{
		DeploymentName: deploymentName,
	}
}

func (id TenantTemplateDeploymentId) String() string {
	segments := []string{
		fmt.Sprintf("Deployment Name %q", id.DeploymentName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Tenant Template Deployment", segmentsStr)
}

func (id TenantTemplateDeploymentId) ID() string {
	fmtString := "/providers/Microsoft.Resources/deployments/%s"
	return fmt.Sprintf(fmtString, id.DeploymentName)
}

// TenantTemplateDeploymentID parses a TenantTemplateDeployment ID into an TenantTemplateDeploymentId struct
func TenantTemplateDeploymentID(input string) (*TenantTemplateDeploymentId, error) {
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

	resourceId := TenantTemplateDeploymentId{}

	if resourceId.DeploymentName, err = id.PopSegment("deployments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
