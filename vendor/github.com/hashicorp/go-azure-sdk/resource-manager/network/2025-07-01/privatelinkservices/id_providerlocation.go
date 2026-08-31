package privatelinkservices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProviderLocationId{})
}

var _ resourceids.ResourceId = &ProviderLocationId{}

// ProviderLocationId is a struct representing the Resource ID for a Provider Location
type ProviderLocationId struct {
	SubscriptionId    string
	ResourceGroupName string
	LocationName      string
}

// NewProviderLocationID returns a new ProviderLocationId struct
func NewProviderLocationID(subscriptionId string, resourceGroupName string, locationName string) ProviderLocationId {
	return ProviderLocationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LocationName:      locationName,
	}
}

// ParseProviderLocationID parses 'input' into a ProviderLocationId
func ParseProviderLocationID(input string) (*ProviderLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderLocationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderLocationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviderLocationIDInsensitively parses 'input' case-insensitively into a ProviderLocationId
// note: this method should only be used for API response data and not user input
func ParseProviderLocationIDInsensitively(input string) (*ProviderLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderLocationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderLocationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProviderLocationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	return nil
}

// ValidateProviderLocationID checks that 'input' can be parsed as a Provider Location ID
func ValidateProviderLocationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderLocationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Location ID
func (id ProviderLocationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/locations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Location ID
func (id ProviderLocationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
	}
}

// String returns a human-readable description of this Provider Location ID
func (id ProviderLocationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Location Name: %q", id.LocationName),
	}
	return fmt.Sprintf("Provider Location (%s)", strings.Join(components, "\n"))
}
