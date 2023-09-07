package partnerregistrations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PartnerRegistrationId{}

// PartnerRegistrationId is a struct representing the Resource ID for a Partner Registration
type PartnerRegistrationId struct {
	SubscriptionId          string
	ResourceGroupName       string
	PartnerRegistrationName string
}

// NewPartnerRegistrationID returns a new PartnerRegistrationId struct
func NewPartnerRegistrationID(subscriptionId string, resourceGroupName string, partnerRegistrationName string) PartnerRegistrationId {
	return PartnerRegistrationId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		PartnerRegistrationName: partnerRegistrationName,
	}
}

// ParsePartnerRegistrationID parses 'input' into a PartnerRegistrationId
func ParsePartnerRegistrationID(input string) (*PartnerRegistrationId, error) {
	parser := resourceids.NewParserFromResourceIdType(PartnerRegistrationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PartnerRegistrationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PartnerRegistrationName, ok = parsed.Parsed["partnerRegistrationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "partnerRegistrationName", *parsed)
	}

	return &id, nil
}

// ParsePartnerRegistrationIDInsensitively parses 'input' case-insensitively into a PartnerRegistrationId
// note: this method should only be used for API response data and not user input
func ParsePartnerRegistrationIDInsensitively(input string) (*PartnerRegistrationId, error) {
	parser := resourceids.NewParserFromResourceIdType(PartnerRegistrationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PartnerRegistrationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PartnerRegistrationName, ok = parsed.Parsed["partnerRegistrationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "partnerRegistrationName", *parsed)
	}

	return &id, nil
}

// ValidatePartnerRegistrationID checks that 'input' can be parsed as a Partner Registration ID
func ValidatePartnerRegistrationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePartnerRegistrationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Partner Registration ID
func (id PartnerRegistrationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/partnerRegistrations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PartnerRegistrationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Partner Registration ID
func (id PartnerRegistrationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticPartnerRegistrations", "partnerRegistrations", "partnerRegistrations"),
		resourceids.UserSpecifiedSegment("partnerRegistrationName", "partnerRegistrationValue"),
	}
}

// String returns a human-readable description of this Partner Registration ID
func (id PartnerRegistrationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Partner Registration Name: %q", id.PartnerRegistrationName),
	}
	return fmt.Sprintf("Partner Registration (%s)", strings.Join(components, "\n"))
}
