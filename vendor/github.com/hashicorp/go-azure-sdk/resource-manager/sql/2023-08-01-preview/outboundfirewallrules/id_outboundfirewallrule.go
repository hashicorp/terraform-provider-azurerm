package outboundfirewallrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OutboundFirewallRuleId{})
}

var _ resourceids.ResourceId = &OutboundFirewallRuleId{}

// OutboundFirewallRuleId is a struct representing the Resource ID for a Outbound Firewall Rule
type OutboundFirewallRuleId struct {
	SubscriptionId           string
	ResourceGroupName        string
	ServerName               string
	OutboundFirewallRuleName string
}

// NewOutboundFirewallRuleID returns a new OutboundFirewallRuleId struct
func NewOutboundFirewallRuleID(subscriptionId string, resourceGroupName string, serverName string, outboundFirewallRuleName string) OutboundFirewallRuleId {
	return OutboundFirewallRuleId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		ServerName:               serverName,
		OutboundFirewallRuleName: outboundFirewallRuleName,
	}
}

// ParseOutboundFirewallRuleID parses 'input' into a OutboundFirewallRuleId
func ParseOutboundFirewallRuleID(input string) (*OutboundFirewallRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OutboundFirewallRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OutboundFirewallRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOutboundFirewallRuleIDInsensitively parses 'input' case-insensitively into a OutboundFirewallRuleId
// note: this method should only be used for API response data and not user input
func ParseOutboundFirewallRuleIDInsensitively(input string) (*OutboundFirewallRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OutboundFirewallRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OutboundFirewallRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OutboundFirewallRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServerName, ok = input.Parsed["serverName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serverName", input)
	}

	if id.OutboundFirewallRuleName, ok = input.Parsed["outboundFirewallRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "outboundFirewallRuleName", input)
	}

	return nil
}

// ValidateOutboundFirewallRuleID checks that 'input' can be parsed as a Outbound Firewall Rule ID
func ValidateOutboundFirewallRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOutboundFirewallRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Outbound Firewall Rule ID
func (id OutboundFirewallRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/outboundFirewallRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.OutboundFirewallRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Outbound Firewall Rule ID
func (id OutboundFirewallRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticServers", "servers", "servers"),
		resourceids.UserSpecifiedSegment("serverName", "serverName"),
		resourceids.StaticSegment("staticOutboundFirewallRules", "outboundFirewallRules", "outboundFirewallRules"),
		resourceids.UserSpecifiedSegment("outboundFirewallRuleName", "outboundFirewallRuleName"),
	}
}

// String returns a human-readable description of this Outbound Firewall Rule ID
func (id OutboundFirewallRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
		fmt.Sprintf("Outbound Firewall Rule Name: %q", id.OutboundFirewallRuleName),
	}
	return fmt.Sprintf("Outbound Firewall Rule (%s)", strings.Join(components, "\n"))
}
