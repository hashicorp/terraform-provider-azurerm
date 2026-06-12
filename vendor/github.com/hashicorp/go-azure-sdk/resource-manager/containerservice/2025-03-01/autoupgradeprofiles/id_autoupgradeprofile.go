package autoupgradeprofiles

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AutoUpgradeProfileId{})
}

var _ resourceids.ResourceId = &AutoUpgradeProfileId{}

// AutoUpgradeProfileId is a struct representing the Resource ID for a Auto Upgrade Profile
type AutoUpgradeProfileId struct {
	SubscriptionId         string
	ResourceGroupName      string
	FleetName              string
	AutoUpgradeProfileName string
}

// NewAutoUpgradeProfileID returns a new AutoUpgradeProfileId struct
func NewAutoUpgradeProfileID(subscriptionId string, resourceGroupName string, fleetName string, autoUpgradeProfileName string) AutoUpgradeProfileId {
	return AutoUpgradeProfileId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		FleetName:              fleetName,
		AutoUpgradeProfileName: autoUpgradeProfileName,
	}
}

// ParseAutoUpgradeProfileID parses 'input' into a AutoUpgradeProfileId
func ParseAutoUpgradeProfileID(input string) (*AutoUpgradeProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutoUpgradeProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutoUpgradeProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutoUpgradeProfileIDInsensitively parses 'input' case-insensitively into a AutoUpgradeProfileId
// note: this method should only be used for API response data and not user input
func ParseAutoUpgradeProfileIDInsensitively(input string) (*AutoUpgradeProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutoUpgradeProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutoUpgradeProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutoUpgradeProfileId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FleetName, ok = input.Parsed["fleetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "fleetName", input)
	}

	if id.AutoUpgradeProfileName, ok = input.Parsed["autoUpgradeProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "autoUpgradeProfileName", input)
	}

	return nil
}

// ValidateAutoUpgradeProfileID checks that 'input' can be parsed as a Auto Upgrade Profile ID
func ValidateAutoUpgradeProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutoUpgradeProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Auto Upgrade Profile ID
func (id AutoUpgradeProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/fleets/%s/autoUpgradeProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FleetName, id.AutoUpgradeProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Auto Upgrade Profile ID
func (id AutoUpgradeProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticFleets", "fleets", "fleets"),
		resourceids.UserSpecifiedSegment("fleetName", "fleetName"),
		resourceids.StaticSegment("staticAutoUpgradeProfiles", "autoUpgradeProfiles", "autoUpgradeProfiles"),
		resourceids.UserSpecifiedSegment("autoUpgradeProfileName", "autoUpgradeProfileName"),
	}
}

// String returns a human-readable description of this Auto Upgrade Profile ID
func (id AutoUpgradeProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Fleet Name: %q", id.FleetName),
		fmt.Sprintf("Auto Upgrade Profile Name: %q", id.AutoUpgradeProfileName),
	}
	return fmt.Sprintf("Auto Upgrade Profile (%s)", strings.Join(components, "\n"))
}
