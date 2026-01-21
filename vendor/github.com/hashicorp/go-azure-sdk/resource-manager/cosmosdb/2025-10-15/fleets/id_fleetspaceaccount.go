package fleets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FleetspaceAccountId{})
}

var _ resourceids.ResourceId = &FleetspaceAccountId{}

// FleetspaceAccountId is a struct representing the Resource ID for a Fleetspace Account
type FleetspaceAccountId struct {
	SubscriptionId        string
	ResourceGroupName     string
	FleetName             string
	FleetspaceName        string
	FleetspaceAccountName string
}

// NewFleetspaceAccountID returns a new FleetspaceAccountId struct
func NewFleetspaceAccountID(subscriptionId string, resourceGroupName string, fleetName string, fleetspaceName string, fleetspaceAccountName string) FleetspaceAccountId {
	return FleetspaceAccountId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		FleetName:             fleetName,
		FleetspaceName:        fleetspaceName,
		FleetspaceAccountName: fleetspaceAccountName,
	}
}

// ParseFleetspaceAccountID parses 'input' into a FleetspaceAccountId
func ParseFleetspaceAccountID(input string) (*FleetspaceAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FleetspaceAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FleetspaceAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFleetspaceAccountIDInsensitively parses 'input' case-insensitively into a FleetspaceAccountId
// note: this method should only be used for API response data and not user input
func ParseFleetspaceAccountIDInsensitively(input string) (*FleetspaceAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FleetspaceAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FleetspaceAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FleetspaceAccountId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.FleetspaceName, ok = input.Parsed["fleetspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "fleetspaceName", input)
	}

	if id.FleetspaceAccountName, ok = input.Parsed["fleetspaceAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "fleetspaceAccountName", input)
	}

	return nil
}

// ValidateFleetspaceAccountID checks that 'input' can be parsed as a Fleetspace Account ID
func ValidateFleetspaceAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFleetspaceAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Fleetspace Account ID
func (id FleetspaceAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/fleets/%s/fleetspaces/%s/fleetspaceAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FleetName, id.FleetspaceName, id.FleetspaceAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Fleetspace Account ID
func (id FleetspaceAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticFleets", "fleets", "fleets"),
		resourceids.UserSpecifiedSegment("fleetName", "fleetName"),
		resourceids.StaticSegment("staticFleetspaces", "fleetspaces", "fleetspaces"),
		resourceids.UserSpecifiedSegment("fleetspaceName", "fleetspaceName"),
		resourceids.StaticSegment("staticFleetspaceAccounts", "fleetspaceAccounts", "fleetspaceAccounts"),
		resourceids.UserSpecifiedSegment("fleetspaceAccountName", "fleetspaceAccountName"),
	}
}

// String returns a human-readable description of this Fleetspace Account ID
func (id FleetspaceAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Fleet Name: %q", id.FleetName),
		fmt.Sprintf("Fleetspace Name: %q", id.FleetspaceName),
		fmt.Sprintf("Fleetspace Account Name: %q", id.FleetspaceAccountName),
	}
	return fmt.Sprintf("Fleetspace Account (%s)", strings.Join(components, "\n"))
}
