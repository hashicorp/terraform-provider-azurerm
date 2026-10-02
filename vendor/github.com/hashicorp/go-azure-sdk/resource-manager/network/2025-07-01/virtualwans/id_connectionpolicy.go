package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConnectionPolicyId{})
}

var _ resourceids.ResourceId = &ConnectionPolicyId{}

// ConnectionPolicyId is a struct representing the Resource ID for a Connection Policy
type ConnectionPolicyId struct {
	SubscriptionId       string
	ResourceGroupName    string
	VirtualHubName       string
	ConnectionPolicyName string
}

// NewConnectionPolicyID returns a new ConnectionPolicyId struct
func NewConnectionPolicyID(subscriptionId string, resourceGroupName string, virtualHubName string, connectionPolicyName string) ConnectionPolicyId {
	return ConnectionPolicyId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		VirtualHubName:       virtualHubName,
		ConnectionPolicyName: connectionPolicyName,
	}
}

// ParseConnectionPolicyID parses 'input' into a ConnectionPolicyId
func ParseConnectionPolicyID(input string) (*ConnectionPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectionPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectionPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConnectionPolicyIDInsensitively parses 'input' case-insensitively into a ConnectionPolicyId
// note: this method should only be used for API response data and not user input
func ParseConnectionPolicyIDInsensitively(input string) (*ConnectionPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectionPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectionPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConnectionPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualHubName, ok = input.Parsed["virtualHubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualHubName", input)
	}

	if id.ConnectionPolicyName, ok = input.Parsed["connectionPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectionPolicyName", input)
	}

	return nil
}

// ValidateConnectionPolicyID checks that 'input' can be parsed as a Connection Policy ID
func ValidateConnectionPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectionPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connection Policy ID
func (id ConnectionPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s/connectionPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, id.ConnectionPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connection Policy ID
func (id ConnectionPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualHubs", "virtualHubs", "virtualHubs"),
		resourceids.UserSpecifiedSegment("virtualHubName", "virtualHubName"),
		resourceids.StaticSegment("staticConnectionPolicies", "connectionPolicies", "connectionPolicies"),
		resourceids.UserSpecifiedSegment("connectionPolicyName", "connectionPolicyName"),
	}
}

// String returns a human-readable description of this Connection Policy ID
func (id ConnectionPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Hub Name: %q", id.VirtualHubName),
		fmt.Sprintf("Connection Policy Name: %q", id.ConnectionPolicyName),
	}
	return fmt.Sprintf("Connection Policy (%s)", strings.Join(components, "\n"))
}
