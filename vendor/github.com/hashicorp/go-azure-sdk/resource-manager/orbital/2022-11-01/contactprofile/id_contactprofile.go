package contactprofile

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ContactProfileId{}

// ContactProfileId is a struct representing the Resource ID for a Contact Profile
type ContactProfileId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ContactProfileName string
}

// NewContactProfileID returns a new ContactProfileId struct
func NewContactProfileID(subscriptionId string, resourceGroupName string, contactProfileName string) ContactProfileId {
	return ContactProfileId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ContactProfileName: contactProfileName,
	}
}

// ParseContactProfileID parses 'input' into a ContactProfileId
func ParseContactProfileID(input string) (*ContactProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(ContactProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ContactProfileId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ContactProfileName, ok = parsed.Parsed["contactProfileName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "contactProfileName", *parsed)
	}

	return &id, nil
}

// ParseContactProfileIDInsensitively parses 'input' case-insensitively into a ContactProfileId
// note: this method should only be used for API response data and not user input
func ParseContactProfileIDInsensitively(input string) (*ContactProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(ContactProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ContactProfileId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ContactProfileName, ok = parsed.Parsed["contactProfileName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "contactProfileName", *parsed)
	}

	return &id, nil
}

// ValidateContactProfileID checks that 'input' can be parsed as a Contact Profile ID
func ValidateContactProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseContactProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Contact Profile ID
func (id ContactProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Orbital/contactProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ContactProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Contact Profile ID
func (id ContactProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOrbital", "Microsoft.Orbital", "Microsoft.Orbital"),
		resourceids.StaticSegment("staticContactProfiles", "contactProfiles", "contactProfiles"),
		resourceids.UserSpecifiedSegment("contactProfileName", "contactProfileValue"),
	}
}

// String returns a human-readable description of this Contact Profile ID
func (id ContactProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Contact Profile Name: %q", id.ContactProfileName),
	}
	return fmt.Sprintf("Contact Profile (%s)", strings.Join(components, "\n"))
}
