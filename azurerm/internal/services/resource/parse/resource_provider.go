package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ResourceProviderId struct {
	SubscriptionId   string
	ResourceProvider string
}

func NewResourceProviderID(subscriptionId, resourceProvider string) ResourceProviderId {
	return ResourceProviderId{
		SubscriptionId:   subscriptionId,
		ResourceProvider: resourceProvider,
	}
}

func (id ResourceProviderId) ID() string {
	fmtString := "/subscriptions/%s/providers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceProvider)
}

// ResourceProviderID parses a ResourceProvider ID into an ResourceProviderId struct
func ResourceProviderID(input string) (*ResourceProviderId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ResourceProviderId{
		SubscriptionId:   id.SubscriptionID,
		ResourceProvider: id.Provider,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceProvider == "" {
		return nil, fmt.Errorf("ID was missing the 'providers' element")
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
