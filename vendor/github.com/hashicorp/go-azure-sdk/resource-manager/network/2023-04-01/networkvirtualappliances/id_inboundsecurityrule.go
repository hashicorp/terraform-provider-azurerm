package networkvirtualappliances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = InboundSecurityRuleId{}

// InboundSecurityRuleId is a struct representing the Resource ID for a Inbound Security Rule
type InboundSecurityRuleId struct {
	SubscriptionId              string
	ResourceGroupName           string
	NetworkVirtualApplianceName string
	InboundSecurityRuleName     string
}

// NewInboundSecurityRuleID returns a new InboundSecurityRuleId struct
func NewInboundSecurityRuleID(subscriptionId string, resourceGroupName string, networkVirtualApplianceName string, inboundSecurityRuleName string) InboundSecurityRuleId {
	return InboundSecurityRuleId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		NetworkVirtualApplianceName: networkVirtualApplianceName,
		InboundSecurityRuleName:     inboundSecurityRuleName,
	}
}

// ParseInboundSecurityRuleID parses 'input' into a InboundSecurityRuleId
func ParseInboundSecurityRuleID(input string) (*InboundSecurityRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(InboundSecurityRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InboundSecurityRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkVirtualApplianceName, ok = parsed.Parsed["networkVirtualApplianceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceName", *parsed)
	}

	if id.InboundSecurityRuleName, ok = parsed.Parsed["inboundSecurityRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "inboundSecurityRuleName", *parsed)
	}

	return &id, nil
}

// ParseInboundSecurityRuleIDInsensitively parses 'input' case-insensitively into a InboundSecurityRuleId
// note: this method should only be used for API response data and not user input
func ParseInboundSecurityRuleIDInsensitively(input string) (*InboundSecurityRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(InboundSecurityRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InboundSecurityRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkVirtualApplianceName, ok = parsed.Parsed["networkVirtualApplianceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceName", *parsed)
	}

	if id.InboundSecurityRuleName, ok = parsed.Parsed["inboundSecurityRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "inboundSecurityRuleName", *parsed)
	}

	return &id, nil
}

// ValidateInboundSecurityRuleID checks that 'input' can be parsed as a Inbound Security Rule ID
func ValidateInboundSecurityRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInboundSecurityRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Inbound Security Rule ID
func (id InboundSecurityRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkVirtualAppliances/%s/inboundSecurityRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkVirtualApplianceName, id.InboundSecurityRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Inbound Security Rule ID
func (id InboundSecurityRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkVirtualAppliances", "networkVirtualAppliances", "networkVirtualAppliances"),
		resourceids.UserSpecifiedSegment("networkVirtualApplianceName", "networkVirtualApplianceValue"),
		resourceids.StaticSegment("staticInboundSecurityRules", "inboundSecurityRules", "inboundSecurityRules"),
		resourceids.UserSpecifiedSegment("inboundSecurityRuleName", "inboundSecurityRuleValue"),
	}
}

// String returns a human-readable description of this Inbound Security Rule ID
func (id InboundSecurityRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Virtual Appliance Name: %q", id.NetworkVirtualApplianceName),
		fmt.Sprintf("Inbound Security Rule Name: %q", id.InboundSecurityRuleName),
	}
	return fmt.Sprintf("Inbound Security Rule (%s)", strings.Join(components, "\n"))
}
