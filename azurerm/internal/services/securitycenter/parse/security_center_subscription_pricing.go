package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecurityCenterSubscriptionPricingId struct {
	SubscriptionId string
	PricingName    string
}

func NewSecurityCenterSubscriptionPricingID(subscriptionId, pricingName string) SecurityCenterSubscriptionPricingId {
	return SecurityCenterSubscriptionPricingId{
		SubscriptionId: subscriptionId,
		PricingName:    pricingName,
	}
}

func (id SecurityCenterSubscriptionPricingId) String() string {
	segments := []string{
		fmt.Sprintf("Pricing Name %q", id.PricingName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Security Center Subscription Pricing", segmentsStr)
}

func (id SecurityCenterSubscriptionPricingId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/pricings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.PricingName)
}

// SecurityCenterSubscriptionPricingID parses a SecurityCenterSubscriptionPricing ID into an SecurityCenterSubscriptionPricingId struct
func SecurityCenterSubscriptionPricingID(input string) (*SecurityCenterSubscriptionPricingId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SecurityCenterSubscriptionPricingId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.PricingName, err = id.PopSegment("pricings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
