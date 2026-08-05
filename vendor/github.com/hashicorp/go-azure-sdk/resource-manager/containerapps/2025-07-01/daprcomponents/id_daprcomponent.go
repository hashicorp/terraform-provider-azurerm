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
	recaser.RegisterResourceId(&DaprComponentId{})
}

var _ resourceids.ResourceId = &DaprComponentId{}

// DaprComponentId is a struct representing the Resource ID for a Dapr Component
type DaprComponentId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ManagedEnvironmentName string
	DaprComponentName      string
}

// NewDaprComponentID returns a new DaprComponentId struct
func NewDaprComponentID(subscriptionId string, resourceGroupName string, managedEnvironmentName string, daprComponentName string) DaprComponentId {
	return DaprComponentId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ManagedEnvironmentName: managedEnvironmentName,
		DaprComponentName:      daprComponentName,
	}
}

// ParseDaprComponentID parses 'input' into a DaprComponentId
func ParseDaprComponentID(input string) (*DaprComponentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DaprComponentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DaprComponentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDaprComponentIDInsensitively parses 'input' case-insensitively into a DaprComponentId
// note: this method should only be used for API response data and not user input
func ParseDaprComponentIDInsensitively(input string) (*DaprComponentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DaprComponentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DaprComponentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DaprComponentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ManagedEnvironmentName, ok = input.Parsed["managedEnvironmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedEnvironmentName", input)
	}

	if id.DaprComponentName, ok = input.Parsed["daprComponentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "daprComponentName", input)
	}

	return nil
}

// ValidateDaprComponentID checks that 'input' can be parsed as a Dapr Component ID
func ValidateDaprComponentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDaprComponentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dapr Component ID
func (id DaprComponentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s/daprComponents/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName, id.DaprComponentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dapr Component ID
func (id DaprComponentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticManagedEnvironments", "managedEnvironments", "managedEnvironments"),
		resourceids.UserSpecifiedSegment("managedEnvironmentName", "managedEnvironmentName"),
		resourceids.StaticSegment("staticDaprComponents", "daprComponents", "daprComponents"),
		resourceids.UserSpecifiedSegment("daprComponentName", "daprComponentName"),
	}
}

// String returns a human-readable description of this Dapr Component ID
func (id DaprComponentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Environment Name: %q", id.ManagedEnvironmentName),
		fmt.Sprintf("Dapr Component Name: %q", id.DaprComponentName),
	}
	return fmt.Sprintf("Dapr Component (%s)", strings.Join(components, "\n"))
}
