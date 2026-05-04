package contact

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ContactId{})
}

var _ resourceids.ResourceId = &ContactId{}

// ContactId is a struct representing the Resource ID for a Contact
type ContactId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpacecraftName    string
	ContactName       string
}

// NewContactID returns a new ContactId struct
func NewContactID(subscriptionId string, resourceGroupName string, spacecraftName string, contactName string) ContactId {
	return ContactId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpacecraftName:    spacecraftName,
		ContactName:       contactName,
	}
}

// ParseContactID parses 'input' into a ContactId
func ParseContactID(input string) (*ContactId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContactId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContactId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseContactIDInsensitively parses 'input' case-insensitively into a ContactId
// note: this method should only be used for API response data and not user input
func ParseContactIDInsensitively(input string) (*ContactId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContactId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContactId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ContactId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpacecraftName, ok = input.Parsed["spacecraftName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "spacecraftName", input)
	}

	if id.ContactName, ok = input.Parsed["contactName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "contactName", input)
	}

	return nil
}

// ValidateContactID checks that 'input' can be parsed as a Contact ID
func ValidateContactID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseContactID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Contact ID
func (id ContactId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Orbital/spacecrafts/%s/contacts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpacecraftName, id.ContactName)
}

// Segments returns a slice of Resource ID Segments which comprise this Contact ID
func (id ContactId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOrbital", "Microsoft.Orbital", "Microsoft.Orbital"),
		resourceids.StaticSegment("staticSpacecrafts", "spacecrafts", "spacecrafts"),
		resourceids.UserSpecifiedSegment("spacecraftName", "spacecraftName"),
		resourceids.StaticSegment("staticContacts", "contacts", "contacts"),
		resourceids.UserSpecifiedSegment("contactName", "contactName"),
	}
}

// String returns a human-readable description of this Contact ID
func (id ContactId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spacecraft Name: %q", id.SpacecraftName),
		fmt.Sprintf("Contact Name: %q", id.ContactName),
	}
	return fmt.Sprintf("Contact (%s)", strings.Join(components, "\n"))
}
