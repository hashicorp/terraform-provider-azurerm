package integrationaccountpartners

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PartnerId{}

// PartnerId is a struct representing the Resource ID for a Partner
type PartnerId struct {
	SubscriptionId         string
	ResourceGroupName      string
	IntegrationAccountName string
	PartnerName            string
}

// NewPartnerID returns a new PartnerId struct
func NewPartnerID(subscriptionId string, resourceGroupName string, integrationAccountName string, partnerName string) PartnerId {
	return PartnerId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		IntegrationAccountName: integrationAccountName,
		PartnerName:            partnerName,
	}
}

// ParsePartnerID parses 'input' into a PartnerId
func ParsePartnerID(input string) (*PartnerId, error) {
	parser := resourceids.NewParserFromResourceIdType(PartnerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PartnerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.IntegrationAccountName, ok = parsed.Parsed["integrationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "integrationAccountName", *parsed)
	}

	if id.PartnerName, ok = parsed.Parsed["partnerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "partnerName", *parsed)
	}

	return &id, nil
}

// ParsePartnerIDInsensitively parses 'input' case-insensitively into a PartnerId
// note: this method should only be used for API response data and not user input
func ParsePartnerIDInsensitively(input string) (*PartnerId, error) {
	parser := resourceids.NewParserFromResourceIdType(PartnerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PartnerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.IntegrationAccountName, ok = parsed.Parsed["integrationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "integrationAccountName", *parsed)
	}

	if id.PartnerName, ok = parsed.Parsed["partnerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "partnerName", *parsed)
	}

	return &id, nil
}

// ValidatePartnerID checks that 'input' can be parsed as a Partner ID
func ValidatePartnerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePartnerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Partner ID
func (id PartnerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/partners/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IntegrationAccountName, id.PartnerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Partner ID
func (id PartnerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticIntegrationAccounts", "integrationAccounts", "integrationAccounts"),
		resourceids.UserSpecifiedSegment("integrationAccountName", "integrationAccountValue"),
		resourceids.StaticSegment("staticPartners", "partners", "partners"),
		resourceids.UserSpecifiedSegment("partnerName", "partnerValue"),
	}
}

// String returns a human-readable description of this Partner ID
func (id PartnerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Integration Account Name: %q", id.IntegrationAccountName),
		fmt.Sprintf("Partner Name: %q", id.PartnerName),
	}
	return fmt.Sprintf("Partner (%s)", strings.Join(components, "\n"))
}
