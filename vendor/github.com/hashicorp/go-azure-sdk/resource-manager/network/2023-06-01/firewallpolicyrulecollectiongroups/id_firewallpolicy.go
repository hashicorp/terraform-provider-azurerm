package firewallpolicyrulecollectiongroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FirewallPolicyId{}

// FirewallPolicyId is a struct representing the Resource ID for a Firewall Policy
type FirewallPolicyId struct {
	SubscriptionId     string
	ResourceGroupName  string
	FirewallPolicyName string
}

// NewFirewallPolicyID returns a new FirewallPolicyId struct
func NewFirewallPolicyID(subscriptionId string, resourceGroupName string, firewallPolicyName string) FirewallPolicyId {
	return FirewallPolicyId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		FirewallPolicyName: firewallPolicyName,
	}
}

// ParseFirewallPolicyID parses 'input' into a FirewallPolicyId
func ParseFirewallPolicyID(input string) (*FirewallPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(FirewallPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FirewallPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FirewallPolicyName, ok = parsed.Parsed["firewallPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "firewallPolicyName", *parsed)
	}

	return &id, nil
}

// ParseFirewallPolicyIDInsensitively parses 'input' case-insensitively into a FirewallPolicyId
// note: this method should only be used for API response data and not user input
func ParseFirewallPolicyIDInsensitively(input string) (*FirewallPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(FirewallPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FirewallPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FirewallPolicyName, ok = parsed.Parsed["firewallPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "firewallPolicyName", *parsed)
	}

	return &id, nil
}

// ValidateFirewallPolicyID checks that 'input' can be parsed as a Firewall Policy ID
func ValidateFirewallPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFirewallPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Firewall Policy ID
func (id FirewallPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/firewallPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FirewallPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Firewall Policy ID
func (id FirewallPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticFirewallPolicies", "firewallPolicies", "firewallPolicies"),
		resourceids.UserSpecifiedSegment("firewallPolicyName", "firewallPolicyValue"),
	}
}

// String returns a human-readable description of this Firewall Policy ID
func (id FirewallPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Firewall Policy Name: %q", id.FirewallPolicyName),
	}
	return fmt.Sprintf("Firewall Policy (%s)", strings.Join(components, "\n"))
}
