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
	recaser.RegisterResourceId(&OfferPlanId{})
}

var _ resourceids.ResourceId = &OfferPlanId{}

// OfferPlanId is a struct representing the Resource ID for a Offer Plan
type OfferPlanId struct {
	SubscriptionId string
	PublisherId    string
	OfferId        string
	PlanId         string
}

// NewOfferPlanID returns a new OfferPlanId struct
func NewOfferPlanID(subscriptionId string, publisherId string, offerId string, planId string) OfferPlanId {
	return OfferPlanId{
		SubscriptionId: subscriptionId,
		PublisherId:    publisherId,
		OfferId:        offerId,
		PlanId:         planId,
	}
}

// ParseOfferPlanID parses 'input' into a OfferPlanId
func ParseOfferPlanID(input string) (*OfferPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OfferPlanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OfferPlanId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOfferPlanIDInsensitively parses 'input' case-insensitively into a OfferPlanId
// note: this method should only be used for API response data and not user input
func ParseOfferPlanIDInsensitively(input string) (*OfferPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OfferPlanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OfferPlanId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OfferPlanId) FromParseResult(input resourceids.ParseResult) error {
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

// ValidateOfferPlanID checks that 'input' can be parsed as a Offer Plan ID
func ValidateOfferPlanID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOfferPlanID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Offer Plan ID
func (id OfferPlanId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.MarketplaceOrdering/offerTypes/virtualMachine/publishers/%s/offers/%s/plans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.PublisherId, id.OfferId, id.PlanId)
}

// Segments returns a slice of Resource ID Segments which comprise this Offer Plan ID
func (id OfferPlanId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMarketplaceOrdering", "Microsoft.MarketplaceOrdering", "Microsoft.MarketplaceOrdering"),
		resourceids.StaticSegment("staticOfferTypes", "offerTypes", "offerTypes"),
		resourceids.StaticSegment("offerType", "virtualMachine", "virtualMachine"),
		resourceids.StaticSegment("staticPublishers", "publishers", "publishers"),
		resourceids.UserSpecifiedSegment("publisherId", "publisherId"),
		resourceids.StaticSegment("staticOffers", "offers", "offers"),
		resourceids.UserSpecifiedSegment("offerId", "offerId"),
		resourceids.StaticSegment("staticPlans", "plans", "plans"),
		resourceids.UserSpecifiedSegment("planId", "planId"),
	}
}

// String returns a human-readable description of this Offer Plan ID
func (id OfferPlanId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Publisher: %q", id.PublisherId),
		fmt.Sprintf("Offer: %q", id.OfferId),
		fmt.Sprintf("Plan: %q", id.PlanId),
	}
	return fmt.Sprintf("Offer Plan (%s)", strings.Join(components, "\n"))
}
