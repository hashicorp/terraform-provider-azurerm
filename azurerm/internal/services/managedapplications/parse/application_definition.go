package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationDefinitionId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewApplicationDefinitionID(subscriptionId, resourceGroup, name string) ApplicationDefinitionId {
	return ApplicationDefinitionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id ApplicationDefinitionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Solutions/applicationDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ApplicationDefinitionID parses a ApplicationDefinition ID into an ApplicationDefinitionId struct
func ApplicationDefinitionID(input string) (*ApplicationDefinitionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApplicationDefinitionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.Name, err = id.PopSegment("applicationDefinitions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
