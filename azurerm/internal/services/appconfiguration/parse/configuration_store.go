package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ConfigurationStoreId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewConfigurationStoreID(subscriptionId, resourceGroup, name string) ConfigurationStoreId {
	return ConfigurationStoreId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id ConfigurationStoreId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppConfiguration/configurationStores/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ConfigurationStoreID parses a ConfigurationStore ID into an ConfigurationStoreId struct
func ConfigurationStoreID(input string) (*ConfigurationStoreId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConfigurationStoreId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.Name, err = id.PopSegment("configurationStores"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
