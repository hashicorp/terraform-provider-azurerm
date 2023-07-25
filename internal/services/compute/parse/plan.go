// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type PlanId struct {
	SubscriptionId string
	AgreementName  string
	OfferName      string
	Name           string
}

func NewPlanID(subscriptionId, agreementName, offerName, name string) PlanId {
	return PlanId{
		SubscriptionId: subscriptionId,
		AgreementName:  agreementName,
		OfferName:      offerName,
		Name:           name,
	}
}

func (id PlanId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Offer Name %q", id.OfferName),
		fmt.Sprintf("Agreement Name %q", id.AgreementName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Plan", segmentsStr)
}

func (id PlanId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.MarketplaceOrdering/agreements/%s/offers/%s/plans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.AgreementName, id.OfferName, id.Name)
}

// PlanID parses a Plan ID into an PlanId struct
func PlanID(input string) (*PlanId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an Plan ID: %+v", input, err)
	}

	resourceId := PlanId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.AgreementName, err = id.PopSegment("agreements"); err != nil {
		return nil, err
	}
	if resourceId.OfferName, err = id.PopSegment("offers"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("plans"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
