package agreements

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PlanId{}

// PlanId is a struct representing the Resource ID for a Plan
type PlanId struct {
	SubscriptionId string
	PublisherId    string
	OfferId        string
	PlanId         string
}

// NewPlanID returns a new PlanId struct
func NewPlanID(subscriptionId string, publisherId string, offerId string, planId string) PlanId {
	return PlanId{
		SubscriptionId: subscriptionId,
		PublisherId:    publisherId,
		OfferId:        offerId,
		PlanId:         planId,
	}
}

// ParsePlanID parses 'input' into a PlanId
func ParsePlanID(input string) (*PlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(PlanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PlanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.PublisherId, ok = parsed.Parsed["publisherId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "publisherId", *parsed)
	}

	if id.OfferId, ok = parsed.Parsed["offerId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "offerId", *parsed)
	}

	if id.PlanId, ok = parsed.Parsed["planId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "planId", *parsed)
	}

	return &id, nil
}

// ParsePlanIDInsensitively parses 'input' case-insensitively into a PlanId
// note: this method should only be used for API response data and not user input
func ParsePlanIDInsensitively(input string) (*PlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(PlanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PlanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.PublisherId, ok = parsed.Parsed["publisherId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "publisherId", *parsed)
	}

	if id.OfferId, ok = parsed.Parsed["offerId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "offerId", *parsed)
	}

	if id.PlanId, ok = parsed.Parsed["planId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "planId", *parsed)
	}

	return &id, nil
}

// ValidatePlanID checks that 'input' can be parsed as a Plan ID
func ValidatePlanID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePlanID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Plan ID
func (id PlanId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.MarketplaceOrdering/agreements/%s/offers/%s/plans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.PublisherId, id.OfferId, id.PlanId)
}

// Segments returns a slice of Resource ID Segments which comprise this Plan ID
func (id PlanId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMarketplaceOrdering", "Microsoft.MarketplaceOrdering", "Microsoft.MarketplaceOrdering"),
		resourceids.StaticSegment("staticAgreements", "agreements", "agreements"),
		resourceids.UserSpecifiedSegment("publisherId", "publisherIdValue"),
		resourceids.StaticSegment("staticOffers", "offers", "offers"),
		resourceids.UserSpecifiedSegment("offerId", "offerIdValue"),
		resourceids.StaticSegment("staticPlans", "plans", "plans"),
		resourceids.UserSpecifiedSegment("planId", "planIdValue"),
	}
}

// String returns a human-readable description of this Plan ID
func (id PlanId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Publisher: %q", id.PublisherId),
		fmt.Sprintf("Offer: %q", id.OfferId),
		fmt.Sprintf("Plan: %q", id.PlanId),
	}
	return fmt.Sprintf("Plan (%s)", strings.Join(components, "\n"))
}
