package expressroutecrossconnectionpeerings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ExpressRouteCrossConnectionId{}

// ExpressRouteCrossConnectionId is a struct representing the Resource ID for a Express Route Cross Connection
type ExpressRouteCrossConnectionId struct {
	SubscriptionId                  string
	ResourceGroupName               string
	ExpressRouteCrossConnectionName string
}

// NewExpressRouteCrossConnectionID returns a new ExpressRouteCrossConnectionId struct
func NewExpressRouteCrossConnectionID(subscriptionId string, resourceGroupName string, expressRouteCrossConnectionName string) ExpressRouteCrossConnectionId {
	return ExpressRouteCrossConnectionId{
		SubscriptionId:                  subscriptionId,
		ResourceGroupName:               resourceGroupName,
		ExpressRouteCrossConnectionName: expressRouteCrossConnectionName,
	}
}

// ParseExpressRouteCrossConnectionID parses 'input' into a ExpressRouteCrossConnectionId
func ParseExpressRouteCrossConnectionID(input string) (*ExpressRouteCrossConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRouteCrossConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRouteCrossConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRouteCrossConnectionName, ok = parsed.Parsed["expressRouteCrossConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCrossConnectionName", *parsed)
	}

	return &id, nil
}

// ParseExpressRouteCrossConnectionIDInsensitively parses 'input' case-insensitively into a ExpressRouteCrossConnectionId
// note: this method should only be used for API response data and not user input
func ParseExpressRouteCrossConnectionIDInsensitively(input string) (*ExpressRouteCrossConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRouteCrossConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRouteCrossConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRouteCrossConnectionName, ok = parsed.Parsed["expressRouteCrossConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCrossConnectionName", *parsed)
	}

	return &id, nil
}

// ValidateExpressRouteCrossConnectionID checks that 'input' can be parsed as a Express Route Cross Connection ID
func ValidateExpressRouteCrossConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRouteCrossConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Cross Connection ID
func (id ExpressRouteCrossConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCrossConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCrossConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Cross Connection ID
func (id ExpressRouteCrossConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteCrossConnections", "expressRouteCrossConnections", "expressRouteCrossConnections"),
		resourceids.UserSpecifiedSegment("expressRouteCrossConnectionName", "expressRouteCrossConnectionValue"),
	}
}

// String returns a human-readable description of this Express Route Cross Connection ID
func (id ExpressRouteCrossConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Cross Connection Name: %q", id.ExpressRouteCrossConnectionName),
	}
	return fmt.Sprintf("Express Route Cross Connection (%s)", strings.Join(components, "\n"))
}
