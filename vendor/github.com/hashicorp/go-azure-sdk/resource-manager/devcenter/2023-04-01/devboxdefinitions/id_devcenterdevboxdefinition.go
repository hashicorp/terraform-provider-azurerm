package devboxdefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DevCenterDevBoxDefinitionId{})
}

var _ resourceids.ResourceId = &DevCenterDevBoxDefinitionId{}

// DevCenterDevBoxDefinitionId is a struct representing the Resource ID for a Dev Center Dev Box Definition
type DevCenterDevBoxDefinitionId struct {
	SubscriptionId       string
	ResourceGroupName    string
	DevCenterName        string
	DevBoxDefinitionName string
}

// NewDevCenterDevBoxDefinitionID returns a new DevCenterDevBoxDefinitionId struct
func NewDevCenterDevBoxDefinitionID(subscriptionId string, resourceGroupName string, devCenterName string, devBoxDefinitionName string) DevCenterDevBoxDefinitionId {
	return DevCenterDevBoxDefinitionId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		DevCenterName:        devCenterName,
		DevBoxDefinitionName: devBoxDefinitionName,
	}
}

// ParseDevCenterDevBoxDefinitionID parses 'input' into a DevCenterDevBoxDefinitionId
func ParseDevCenterDevBoxDefinitionID(input string) (*DevCenterDevBoxDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DevCenterDevBoxDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DevCenterDevBoxDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDevCenterDevBoxDefinitionIDInsensitively parses 'input' case-insensitively into a DevCenterDevBoxDefinitionId
// note: this method should only be used for API response data and not user input
func ParseDevCenterDevBoxDefinitionIDInsensitively(input string) (*DevCenterDevBoxDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DevCenterDevBoxDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DevCenterDevBoxDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DevCenterDevBoxDefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DevCenterName, ok = input.Parsed["devCenterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "devCenterName", input)
	}

	if id.DevBoxDefinitionName, ok = input.Parsed["devBoxDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "devBoxDefinitionName", input)
	}

	return nil
}

// ValidateDevCenterDevBoxDefinitionID checks that 'input' can be parsed as a Dev Center Dev Box Definition ID
func ValidateDevCenterDevBoxDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDevCenterDevBoxDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dev Center Dev Box Definition ID
func (id DevCenterDevBoxDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/devCenters/%s/devBoxDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DevCenterName, id.DevBoxDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dev Center Dev Box Definition ID
func (id DevCenterDevBoxDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticDevCenters", "devCenters", "devCenters"),
		resourceids.UserSpecifiedSegment("devCenterName", "devCenterName"),
		resourceids.StaticSegment("staticDevBoxDefinitions", "devBoxDefinitions", "devBoxDefinitions"),
		resourceids.UserSpecifiedSegment("devBoxDefinitionName", "devBoxDefinitionName"),
	}
}

// String returns a human-readable description of this Dev Center Dev Box Definition ID
func (id DevCenterDevBoxDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dev Center Name: %q", id.DevCenterName),
		fmt.Sprintf("Dev Box Definition Name: %q", id.DevBoxDefinitionName),
	}
	return fmt.Sprintf("Dev Center Dev Box Definition (%s)", strings.Join(components, "\n"))
}
