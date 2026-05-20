package customlocations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CustomLocationId{})
}

var _ resourceids.ResourceId = &CustomLocationId{}

// CustomLocationId is a struct representing the Resource ID for a Custom Location
type CustomLocationId struct {
	SubscriptionId     string
	ResourceGroupName  string
	CustomLocationName string
}

// NewCustomLocationID returns a new CustomLocationId struct
func NewCustomLocationID(subscriptionId string, resourceGroupName string, customLocationName string) CustomLocationId {
	return CustomLocationId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		CustomLocationName: customLocationName,
	}
}

// ParseCustomLocationID parses 'input' into a CustomLocationId
func ParseCustomLocationID(input string) (*CustomLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CustomLocationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CustomLocationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCustomLocationIDInsensitively parses 'input' case-insensitively into a CustomLocationId
// note: this method should only be used for API response data and not user input
func ParseCustomLocationIDInsensitively(input string) (*CustomLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CustomLocationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CustomLocationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CustomLocationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CustomLocationName, ok = input.Parsed["customLocationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "customLocationName", input)
	}

	return nil
}

// ValidateCustomLocationID checks that 'input' can be parsed as a Custom Location ID
func ValidateCustomLocationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCustomLocationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Custom Location ID
func (id CustomLocationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ExtendedLocation/customLocations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CustomLocationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Custom Location ID
func (id CustomLocationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftExtendedLocation", "Microsoft.ExtendedLocation", "Microsoft.ExtendedLocation"),
		resourceids.StaticSegment("staticCustomLocations", "customLocations", "customLocations"),
		resourceids.UserSpecifiedSegment("customLocationName", "customLocationName"),
	}
}

// String returns a human-readable description of this Custom Location ID
func (id CustomLocationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Custom Location Name: %q", id.CustomLocationName),
	}
	return fmt.Sprintf("Custom Location (%s)", strings.Join(components, "\n"))
}
