package scopeconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopeConnectionId{}

// ScopeConnectionId is a struct representing the Resource ID for a Scope Connection
type ScopeConnectionId struct {
	SubscriptionId      string
	ResourceGroupName   string
	NetworkManagerName  string
	ScopeConnectionName string
}

// NewScopeConnectionID returns a new ScopeConnectionId struct
func NewScopeConnectionID(subscriptionId string, resourceGroupName string, networkManagerName string, scopeConnectionName string) ScopeConnectionId {
	return ScopeConnectionId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		NetworkManagerName:  networkManagerName,
		ScopeConnectionName: scopeConnectionName,
	}
}

// ParseScopeConnectionID parses 'input' into a ScopeConnectionId
func ParseScopeConnectionID(input string) (*ScopeConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopeConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopeConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkManagerName, ok = parsed.Parsed["networkManagerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", *parsed)
	}

	if id.ScopeConnectionName, ok = parsed.Parsed["scopeConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scopeConnectionName", *parsed)
	}

	return &id, nil
}

// ParseScopeConnectionIDInsensitively parses 'input' case-insensitively into a ScopeConnectionId
// note: this method should only be used for API response data and not user input
func ParseScopeConnectionIDInsensitively(input string) (*ScopeConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopeConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopeConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkManagerName, ok = parsed.Parsed["networkManagerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", *parsed)
	}

	if id.ScopeConnectionName, ok = parsed.Parsed["scopeConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scopeConnectionName", *parsed)
	}

	return &id, nil
}

// ValidateScopeConnectionID checks that 'input' can be parsed as a Scope Connection ID
func ValidateScopeConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopeConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scope Connection ID
func (id ScopeConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/scopeConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.ScopeConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scope Connection ID
func (id ScopeConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerValue"),
		resourceids.StaticSegment("staticScopeConnections", "scopeConnections", "scopeConnections"),
		resourceids.UserSpecifiedSegment("scopeConnectionName", "scopeConnectionValue"),
	}
}

// String returns a human-readable description of this Scope Connection ID
func (id ScopeConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
		fmt.Sprintf("Scope Connection Name: %q", id.ScopeConnectionName),
	}
	return fmt.Sprintf("Scope Connection (%s)", strings.Join(components, "\n"))
}
