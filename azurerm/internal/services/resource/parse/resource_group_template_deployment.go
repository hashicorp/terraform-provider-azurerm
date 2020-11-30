package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ResourceGroupTemplateDeploymentId struct {
	ResourceGroup  string
	DeploymentName string
}

func NewResourceGroupTemplateDeploymentID(resourceGroup, name string) ResourceGroupTemplateDeploymentId {
	return ResourceGroupTemplateDeploymentId{
		ResourceGroup:  resourceGroup,
		DeploymentName: name,
	}
}

func ResourceGroupTemplateDeploymentID(input string) (*ResourceGroupTemplateDeploymentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	deploymentId := ResourceGroupTemplateDeploymentId{
		ResourceGroup: id.ResourceGroup,
	}

	if deploymentId.DeploymentName, err = id.PopSegment("deployments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &deploymentId, nil
}

func (id ResourceGroupTemplateDeploymentId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/deployments/%s", subscriptionId, id.ResourceGroup, id.DeploymentName)
}
