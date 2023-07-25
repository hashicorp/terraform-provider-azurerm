package virtualnetworkrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VirtualNetworkRuleId{}

// VirtualNetworkRuleId is a struct representing the Resource ID for a Virtual Network Rule
type VirtualNetworkRuleId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ServerName             string
	VirtualNetworkRuleName string
}

// NewVirtualNetworkRuleID returns a new VirtualNetworkRuleId struct
func NewVirtualNetworkRuleID(subscriptionId string, resourceGroupName string, serverName string, virtualNetworkRuleName string) VirtualNetworkRuleId {
	return VirtualNetworkRuleId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ServerName:             serverName,
		VirtualNetworkRuleName: virtualNetworkRuleName,
	}
}

// ParseVirtualNetworkRuleID parses 'input' into a VirtualNetworkRuleId
func ParseVirtualNetworkRuleID(input string) (*VirtualNetworkRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualNetworkRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualNetworkRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServerName, ok = parsed.Parsed["serverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serverName", *parsed)
	}

	if id.VirtualNetworkRuleName, ok = parsed.Parsed["virtualNetworkRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkRuleName", *parsed)
	}

	return &id, nil
}

// ParseVirtualNetworkRuleIDInsensitively parses 'input' case-insensitively into a VirtualNetworkRuleId
// note: this method should only be used for API response data and not user input
func ParseVirtualNetworkRuleIDInsensitively(input string) (*VirtualNetworkRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualNetworkRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualNetworkRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServerName, ok = parsed.Parsed["serverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serverName", *parsed)
	}

	if id.VirtualNetworkRuleName, ok = parsed.Parsed["virtualNetworkRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkRuleName", *parsed)
	}

	return &id, nil
}

// ValidateVirtualNetworkRuleID checks that 'input' can be parsed as a Virtual Network Rule ID
func ValidateVirtualNetworkRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualNetworkRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Network Rule ID
func (id VirtualNetworkRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMariaDB/servers/%s/virtualNetworkRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.VirtualNetworkRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Network Rule ID
func (id VirtualNetworkRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforMariaDB", "Microsoft.DBforMariaDB", "Microsoft.DBforMariaDB"),
		resourceids.StaticSegment("staticServers", "servers", "servers"),
		resourceids.UserSpecifiedSegment("serverName", "serverValue"),
		resourceids.StaticSegment("staticVirtualNetworkRules", "virtualNetworkRules", "virtualNetworkRules"),
		resourceids.UserSpecifiedSegment("virtualNetworkRuleName", "virtualNetworkRuleValue"),
	}
}

// String returns a human-readable description of this Virtual Network Rule ID
func (id VirtualNetworkRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
		fmt.Sprintf("Virtual Network Rule Name: %q", id.VirtualNetworkRuleName),
	}
	return fmt.Sprintf("Virtual Network Rule (%s)", strings.Join(components, "\n"))
}
