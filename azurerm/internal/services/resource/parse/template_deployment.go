package parse

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TemplateDeploymentId struct {
	ResourceGroup string
	Name string
}

func NewTemplateDeploymentId(resourceGroup, name string) TemplateDeploymentId {
	return TemplateDeploymentId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id TemplateDeploymentId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/deployments/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func TemplateDeploymentID(input string) (*TemplateDeploymentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Availability Set ID %q: %+v", input, err)
	}

	set := TemplateDeploymentId{
		ResourceGroup: id.ResourceGroup,
	}

	// in some circumstances, the ID may have a `Deployments` segment with capitalized D, therefore here to parse it, we have to switch its casing
	if deployment, ok := id.Path["Deployments"]; ok {
		delete(id.Path, "Deployments")
		id.Path["deployments"] = deployment
	}
	if set.Name, err = id.PopSegment("deployments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &set, nil
}
