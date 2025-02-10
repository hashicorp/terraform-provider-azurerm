package managedapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ManagedApiId{})
}

var _ resourceids.ResourceId = &ManagedApiId{}

// ManagedApiId is a struct representing the Resource ID for a Managed Api
type ManagedApiId struct {
	SubscriptionId string
	LocationName   string
	ManagedApiName string
}

// NewManagedApiID returns a new ManagedApiId struct
func NewManagedApiID(subscriptionId string, locationName string, managedApiName string) ManagedApiId {
	return ManagedApiId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		ManagedApiName: managedApiName,
	}
}

// ParseManagedApiID parses 'input' into a ManagedApiId
func ParseManagedApiID(input string) (*ManagedApiId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedApiId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedApiId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseManagedApiIDInsensitively parses 'input' case-insensitively into a ManagedApiId
// note: this method should only be used for API response data and not user input
func ParseManagedApiIDInsensitively(input string) (*ManagedApiId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedApiId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedApiId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ManagedApiId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.ManagedApiName, ok = input.Parsed["managedApiName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedApiName", input)
	}

	return nil
}

// ValidateManagedApiID checks that 'input' can be parsed as a Managed Api ID
func ValidateManagedApiID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedApiID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Api ID
func (id ManagedApiId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Web/locations/%s/managedApis/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.ManagedApiName)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Api ID
func (id ManagedApiId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticManagedApis", "managedApis", "managedApis"),
		resourceids.UserSpecifiedSegment("managedApiName", "managedApiName"),
	}
}

// String returns a human-readable description of this Managed Api ID
func (id ManagedApiId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Managed Api Name: %q", id.ManagedApiName),
	}
	return fmt.Sprintf("Managed Api (%s)", strings.Join(components, "\n"))
}
