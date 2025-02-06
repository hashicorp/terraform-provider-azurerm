package agreements

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PlanId{})
}

var _ resourceids.ResourceId = &PlanId{}

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
	parser := resourceids.NewParserFromResourceIdType(&PlanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PlanId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePlanIDInsensitively parses 'input' case-insensitively into a PlanId
// note: this method should only be used for API response data and not user input
func ParsePlanIDInsensitively(input string) (*PlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PlanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PlanId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PlanId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.PublisherId, ok = input.Parsed["publisherId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "publisherId", input)
	}

	if id.OfferId, ok = input.Parsed["offerId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "offerId", input)
	}

	if id.PlanId, ok = input.Parsed["planId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "planId", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("publisherId", "publisherId"),
		resourceids.StaticSegment("staticOffers", "offers", "offers"),
		resourceids.UserSpecifiedSegment("offerId", "offerId"),
		resourceids.StaticSegment("staticPlans", "plans", "plans"),
		resourceids.UserSpecifiedSegment("planId", "planId"),
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
