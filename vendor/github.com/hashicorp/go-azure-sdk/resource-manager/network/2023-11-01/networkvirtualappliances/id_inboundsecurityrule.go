package networkvirtualappliances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&InboundSecurityRuleId{})
}

var _ resourceids.ResourceId = &InboundSecurityRuleId{}

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
	parser := resourceids.NewParserFromResourceIdType(&InboundSecurityRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InboundSecurityRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseInboundSecurityRuleIDInsensitively parses 'input' case-insensitively into a InboundSecurityRuleId
// note: this method should only be used for API response data and not user input
func ParseInboundSecurityRuleIDInsensitively(input string) (*InboundSecurityRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InboundSecurityRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InboundSecurityRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *InboundSecurityRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkVirtualApplianceName, ok = input.Parsed["networkVirtualApplianceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceName", input)
	}

	if id.InboundSecurityRuleName, ok = input.Parsed["inboundSecurityRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "inboundSecurityRuleName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("networkVirtualApplianceName", "networkVirtualApplianceName"),
		resourceids.StaticSegment("staticInboundSecurityRules", "inboundSecurityRules", "inboundSecurityRules"),
		resourceids.UserSpecifiedSegment("inboundSecurityRuleName", "inboundSecurityRuleName"),
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
