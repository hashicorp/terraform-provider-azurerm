package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type WebTestId struct {
	SubscriptionId string
	ResourceGroup  string
	WebtestName    string
}

func NewWebTestID(subscriptionId, resourceGroup, webtestName string) WebTestId {
	return WebTestId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WebtestName:    webtestName,
	}
}

func (id WebTestId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/microsoft.insights/webtests/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WebtestName)
}

func WebTestID(input string) (*WebTestId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WebTestId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.WebtestName, err = id.PopSegment("webtests"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
