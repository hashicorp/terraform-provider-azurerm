package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type IntegrationRuntimeId struct {
	SubscriptionId string
	ResourceGroup  string
	FactoryName    string
	Name           string
}

func NewIntegrationRuntimeID(subscriptionId, resourceGroup, factoryName, name string) IntegrationRuntimeId {
	return IntegrationRuntimeId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FactoryName:    factoryName,
		Name:           name,
	}
}

func (id IntegrationRuntimeId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s/integrationruntimes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FactoryName, id.Name)
}

func IntegrationRuntimeID(input string) (*IntegrationRuntimeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IntegrationRuntimeId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.FactoryName, err = id.PopSegment("factories"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("integrationruntimes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
