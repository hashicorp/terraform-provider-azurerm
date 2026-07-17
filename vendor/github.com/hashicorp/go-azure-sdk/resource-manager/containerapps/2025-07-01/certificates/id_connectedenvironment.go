package certificates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConnectedEnvironmentId{})
}

var _ resourceids.ResourceId = &ConnectedEnvironmentId{}

// ConnectedEnvironmentId is a struct representing the Resource ID for a Connected Environment
type ConnectedEnvironmentId struct {
	SubscriptionId           string
	ResourceGroupName        string
	ConnectedEnvironmentName string
}

// NewConnectedEnvironmentID returns a new ConnectedEnvironmentId struct
func NewConnectedEnvironmentID(subscriptionId string, resourceGroupName string, connectedEnvironmentName string) ConnectedEnvironmentId {
	return ConnectedEnvironmentId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		ConnectedEnvironmentName: connectedEnvironmentName,
	}
}

// ParseConnectedEnvironmentID parses 'input' into a ConnectedEnvironmentId
func ParseConnectedEnvironmentID(input string) (*ConnectedEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectedEnvironmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectedEnvironmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConnectedEnvironmentIDInsensitively parses 'input' case-insensitively into a ConnectedEnvironmentId
// note: this method should only be used for API response data and not user input
func ParseConnectedEnvironmentIDInsensitively(input string) (*ConnectedEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectedEnvironmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectedEnvironmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConnectedEnvironmentId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateConnectedEnvironmentID checks that 'input' can be parsed as a Connected Environment ID
func ValidateConnectedEnvironmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectedEnvironmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connected Environment ID
func (id ConnectedEnvironmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/connectedEnvironments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ConnectedEnvironmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connected Environment ID
func (id ConnectedEnvironmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticConnectedEnvironments", "connectedEnvironments", "connectedEnvironments"),
		resourceids.UserSpecifiedSegment("connectedEnvironmentName", "connectedEnvironmentName"),
	}
}

// String returns a human-readable description of this Connected Environment ID
func (id ConnectedEnvironmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Connected Environment Name: %q", id.ConnectedEnvironmentName),
	}
	return fmt.Sprintf("Connected Environment (%s)", strings.Join(components, "\n"))
}
