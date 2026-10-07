package integrationaccountmaps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MapId{})
}

var _ resourceids.ResourceId = &MapId{}

// MapId is a struct representing the Resource ID for a Map
type MapId struct {
	SubscriptionId         string
	ResourceGroupName      string
	IntegrationAccountName string
	MapName                string
}

// NewMapID returns a new MapId struct
func NewMapID(subscriptionId string, resourceGroupName string, integrationAccountName string, mapName string) MapId {
	return MapId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		IntegrationAccountName: integrationAccountName,
		MapName:                mapName,
	}
}

// ParseMapID parses 'input' into a MapId
func ParseMapID(input string) (*MapId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MapId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MapId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMapIDInsensitively parses 'input' case-insensitively into a MapId
// note: this method should only be used for API response data and not user input
func ParseMapIDInsensitively(input string) (*MapId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MapId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MapId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MapId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.IntegrationAccountName, ok = input.Parsed["integrationAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "integrationAccountName", input)
	}

	if id.MapName, ok = input.Parsed["mapName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "mapName", input)
	}

	return nil
}

// ValidateMapID checks that 'input' can be parsed as a Map ID
func ValidateMapID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMapID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Map ID
func (id MapId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/maps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IntegrationAccountName, id.MapName)
}

// Segments returns a slice of Resource ID Segments which comprise this Map ID
func (id MapId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticIntegrationAccounts", "integrationAccounts", "integrationAccounts"),
		resourceids.UserSpecifiedSegment("integrationAccountName", "integrationAccountName"),
		resourceids.StaticSegment("staticMaps", "maps", "maps"),
		resourceids.UserSpecifiedSegment("mapName", "mapName"),
	}
}

// String returns a human-readable description of this Map ID
func (id MapId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Integration Account Name: %q", id.IntegrationAccountName),
		fmt.Sprintf("Map Name: %q", id.MapName),
	}
	return fmt.Sprintf("Map (%s)", strings.Join(components, "\n"))
}
