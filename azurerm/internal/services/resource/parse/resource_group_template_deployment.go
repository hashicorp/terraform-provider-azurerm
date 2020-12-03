package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ResourceGroupTemplateDeploymentId struct {
	SubscriptionId string
	ResourceGroup  string
	DeploymentName string
}

func NewResourceGroupTemplateDeploymentID(subscriptionId, resourceGroup, deploymentName string) ResourceGroupTemplateDeploymentId {
	return ResourceGroupTemplateDeploymentId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DeploymentName: deploymentName,
	}
}

func (id ResourceGroupTemplateDeploymentId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/deployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DeploymentName)
}

// ResourceGroupTemplateDeploymentID parses a ResourceGroupTemplateDeployment ID into an ResourceGroupTemplateDeploymentId struct
func ResourceGroupTemplateDeploymentID(input string) (*ResourceGroupTemplateDeploymentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ResourceGroupTemplateDeploymentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DeploymentName, err = id.PopSegment("deployments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
