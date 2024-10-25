package forwardingrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ForwardingRuleId{})
}

var _ resourceids.ResourceId = &ForwardingRuleId{}

// ForwardingRuleId is a struct representing the Resource ID for a Forwarding Rule
type ForwardingRuleId struct {
	SubscriptionId           string
	ResourceGroupName        string
	DnsForwardingRulesetName string
	ForwardingRuleName       string
}

// NewForwardingRuleID returns a new ForwardingRuleId struct
func NewForwardingRuleID(subscriptionId string, resourceGroupName string, dnsForwardingRulesetName string, forwardingRuleName string) ForwardingRuleId {
	return ForwardingRuleId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		DnsForwardingRulesetName: dnsForwardingRulesetName,
		ForwardingRuleName:       forwardingRuleName,
	}
}

// ParseForwardingRuleID parses 'input' into a ForwardingRuleId
func ParseForwardingRuleID(input string) (*ForwardingRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ForwardingRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ForwardingRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseForwardingRuleIDInsensitively parses 'input' case-insensitively into a ForwardingRuleId
// note: this method should only be used for API response data and not user input
func ParseForwardingRuleIDInsensitively(input string) (*ForwardingRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ForwardingRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ForwardingRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ForwardingRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DnsForwardingRulesetName, ok = input.Parsed["dnsForwardingRulesetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dnsForwardingRulesetName", input)
	}

	if id.ForwardingRuleName, ok = input.Parsed["forwardingRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "forwardingRuleName", input)
	}

	return nil
}

// ValidateForwardingRuleID checks that 'input' can be parsed as a Forwarding Rule ID
func ValidateForwardingRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseForwardingRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Forwarding Rule ID
func (id ForwardingRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsForwardingRulesets/%s/forwardingRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DnsForwardingRulesetName, id.ForwardingRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Forwarding Rule ID
func (id ForwardingRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDnsForwardingRulesets", "dnsForwardingRulesets", "dnsForwardingRulesets"),
		resourceids.UserSpecifiedSegment("dnsForwardingRulesetName", "dnsForwardingRulesetName"),
		resourceids.StaticSegment("staticForwardingRules", "forwardingRules", "forwardingRules"),
		resourceids.UserSpecifiedSegment("forwardingRuleName", "forwardingRuleName"),
	}
}

// String returns a human-readable description of this Forwarding Rule ID
func (id ForwardingRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dns Forwarding Ruleset Name: %q", id.DnsForwardingRulesetName),
		fmt.Sprintf("Forwarding Rule Name: %q", id.ForwardingRuleName),
	}
	return fmt.Sprintf("Forwarding Rule (%s)", strings.Join(components, "\n"))
}
