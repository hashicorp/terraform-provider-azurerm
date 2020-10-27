package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SubscriptionTemplateDeploymentId struct {
	Name string
}

func NewSubscriptionTemplateDeploymentID(name string) SubscriptionTemplateDeploymentId {
	return SubscriptionTemplateDeploymentId{
		Name: name,
	}
}

func SubscriptionTemplateDeploymentID(input string) (*SubscriptionTemplateDeploymentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	deploymentId := SubscriptionTemplateDeploymentId{}

	if deploymentId.Name, err = id.PopSegment("deployments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &deploymentId, nil
}

func (id SubscriptionTemplateDeploymentId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Resources/deployments/%s", subscriptionId, id.Name)
}
