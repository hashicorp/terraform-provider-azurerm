package firewallrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FirewallRuleId{})
}

var _ resourceids.ResourceId = &FirewallRuleId{}

// FirewallRuleId is a struct representing the Resource ID for a Firewall Rule
type FirewallRuleId struct {
	SubscriptionId     string
	ResourceGroupName  string
	FlexibleServerName string
	FirewallRuleName   string
}

// NewFirewallRuleID returns a new FirewallRuleId struct
func NewFirewallRuleID(subscriptionId string, resourceGroupName string, flexibleServerName string, firewallRuleName string) FirewallRuleId {
	return FirewallRuleId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		FlexibleServerName: flexibleServerName,
		FirewallRuleName:   firewallRuleName,
	}
}

// ParseFirewallRuleID parses 'input' into a FirewallRuleId
func ParseFirewallRuleID(input string) (*FirewallRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FirewallRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FirewallRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFirewallRuleIDInsensitively parses 'input' case-insensitively into a FirewallRuleId
// note: this method should only be used for API response data and not user input
func ParseFirewallRuleIDInsensitively(input string) (*FirewallRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FirewallRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FirewallRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FirewallRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FlexibleServerName, ok = input.Parsed["flexibleServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "flexibleServerName", input)
	}

	if id.FirewallRuleName, ok = input.Parsed["firewallRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "firewallRuleName", input)
	}

	return nil
}

// ValidateFirewallRuleID checks that 'input' can be parsed as a Firewall Rule ID
func ValidateFirewallRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFirewallRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Firewall Rule ID
func (id FirewallRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMySQL/flexibleServers/%s/firewallRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName, id.FirewallRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Firewall Rule ID
func (id FirewallRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforMySQL", "Microsoft.DBforMySQL", "Microsoft.DBforMySQL"),
		resourceids.StaticSegment("staticFlexibleServers", "flexibleServers", "flexibleServers"),
		resourceids.UserSpecifiedSegment("flexibleServerName", "flexibleServerName"),
		resourceids.StaticSegment("staticFirewallRules", "firewallRules", "firewallRules"),
		resourceids.UserSpecifiedSegment("firewallRuleName", "firewallRuleName"),
	}
}

// String returns a human-readable description of this Firewall Rule ID
func (id FirewallRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Flexible Server Name: %q", id.FlexibleServerName),
		fmt.Sprintf("Firewall Rule Name: %q", id.FirewallRuleName),
	}
	return fmt.Sprintf("Firewall Rule (%s)", strings.Join(components, "\n"))
}
