package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecurityCenterSubscriptionPricingId struct {
	ResourceType string
}

func SecurityCenterSubscriptionPricingID(input string) (*SecurityCenterSubscriptionPricingId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Security Center Subscription Pricing ID %q: %+v", input, err)
	}

	pricing := SecurityCenterSubscriptionPricingId{}

	if pricing.ResourceType, err = id.PopSegment("pricings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &pricing, nil
}
