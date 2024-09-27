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
	recaser.RegisterResourceId(&SpacecraftId{})
}

var _ resourceids.ResourceId = &SpacecraftId{}

// SpacecraftId is a struct representing the Resource ID for a Spacecraft
type SpacecraftId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpacecraftName    string
}

// NewSpacecraftID returns a new SpacecraftId struct
func NewSpacecraftID(subscriptionId string, resourceGroupName string, spacecraftName string) SpacecraftId {
	return SpacecraftId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpacecraftName:    spacecraftName,
	}
}

// ParseSpacecraftID parses 'input' into a SpacecraftId
func ParseSpacecraftID(input string) (*SpacecraftId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SpacecraftId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SpacecraftId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSpacecraftIDInsensitively parses 'input' case-insensitively into a SpacecraftId
// note: this method should only be used for API response data and not user input
func ParseSpacecraftIDInsensitively(input string) (*SpacecraftId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SpacecraftId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SpacecraftId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SpacecraftId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateSpacecraftID checks that 'input' can be parsed as a Spacecraft ID
func ValidateSpacecraftID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSpacecraftID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Spacecraft ID
func (id SpacecraftId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Orbital/spacecrafts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpacecraftName)
}

// Segments returns a slice of Resource ID Segments which comprise this Spacecraft ID
func (id SpacecraftId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOrbital", "Microsoft.Orbital", "Microsoft.Orbital"),
		resourceids.StaticSegment("staticSpacecrafts", "spacecrafts", "spacecrafts"),
		resourceids.UserSpecifiedSegment("spacecraftName", "spacecraftName"),
	}
}

// String returns a human-readable description of this Spacecraft ID
func (id SpacecraftId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spacecraft Name: %q", id.SpacecraftName),
	}
	return fmt.Sprintf("Spacecraft (%s)", strings.Join(components, "\n"))
}
