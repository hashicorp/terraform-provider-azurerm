package daprcomponents

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConnectedEnvironmentDaprComponentId{})
}

var _ resourceids.ResourceId = &ConnectedEnvironmentDaprComponentId{}

// ConnectedEnvironmentDaprComponentId is a struct representing the Resource ID for a Connected Environment Dapr Component
type ConnectedEnvironmentDaprComponentId struct {
	SubscriptionId           string
	ResourceGroupName        string
	ConnectedEnvironmentName string
	DaprComponentName        string
}

// NewConnectedEnvironmentDaprComponentID returns a new ConnectedEnvironmentDaprComponentId struct
func NewConnectedEnvironmentDaprComponentID(subscriptionId string, resourceGroupName string, connectedEnvironmentName string, daprComponentName string) ConnectedEnvironmentDaprComponentId {
	return ConnectedEnvironmentDaprComponentId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		ConnectedEnvironmentName: connectedEnvironmentName,
		DaprComponentName:        daprComponentName,
	}
}

// ParseConnectedEnvironmentDaprComponentID parses 'input' into a ConnectedEnvironmentDaprComponentId
func ParseConnectedEnvironmentDaprComponentID(input string) (*ConnectedEnvironmentDaprComponentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectedEnvironmentDaprComponentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectedEnvironmentDaprComponentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConnectedEnvironmentDaprComponentIDInsensitively parses 'input' case-insensitively into a ConnectedEnvironmentDaprComponentId
// note: this method should only be used for API response data and not user input
func ParseConnectedEnvironmentDaprComponentIDInsensitively(input string) (*ConnectedEnvironmentDaprComponentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectedEnvironmentDaprComponentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectedEnvironmentDaprComponentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConnectedEnvironmentDaprComponentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ConnectedEnvironmentName, ok = input.Parsed["connectedEnvironmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectedEnvironmentName", input)
	}

	if id.DaprComponentName, ok = input.Parsed["daprComponentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "daprComponentName", input)
	}

	return nil
}

// ValidateConnectedEnvironmentDaprComponentID checks that 'input' can be parsed as a Connected Environment Dapr Component ID
func ValidateConnectedEnvironmentDaprComponentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectedEnvironmentDaprComponentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connected Environment Dapr Component ID
func (id ConnectedEnvironmentDaprComponentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/connectedEnvironments/%s/daprComponents/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ConnectedEnvironmentName, id.DaprComponentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connected Environment Dapr Component ID
func (id ConnectedEnvironmentDaprComponentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticConnectedEnvironments", "connectedEnvironments", "connectedEnvironments"),
		resourceids.UserSpecifiedSegment("connectedEnvironmentName", "connectedEnvironmentName"),
		resourceids.StaticSegment("staticDaprComponents", "daprComponents", "daprComponents"),
		resourceids.UserSpecifiedSegment("daprComponentName", "daprComponentName"),
	}
}

// String returns a human-readable description of this Connected Environment Dapr Component ID
func (id ConnectedEnvironmentDaprComponentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Connected Environment Name: %q", id.ConnectedEnvironmentName),
		fmt.Sprintf("Dapr Component Name: %q", id.DaprComponentName),
	}
	return fmt.Sprintf("Connected Environment Dapr Component (%s)", strings.Join(components, "\n"))
}
