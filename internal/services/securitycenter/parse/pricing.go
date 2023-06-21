// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type PricingId struct {
	SubscriptionId string
	Name           string
}

func NewPricingID(subscriptionId, name string) PricingId {
	return PricingId{
		SubscriptionId: subscriptionId,
		Name:           name,
	}
}

func (id PricingId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Pricing", segmentsStr)
}

func (id PricingId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/pricings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Name)
}

// PricingID parses a Pricing ID into an PricingId struct
func PricingID(input string) (*PricingId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an Pricing ID: %+v", input, err)
	}

	resourceId := PricingId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.Name, err = id.PopSegment("pricings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
