package webapplicationfirewallpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ApplicationGatewayWebApplicationFirewallPolicyId{}

// ApplicationGatewayWebApplicationFirewallPolicyId is a struct representing the Resource ID for a Application Gateway Web Application Firewall Policy
type ApplicationGatewayWebApplicationFirewallPolicyId struct {
	SubscriptionId                                     string
	ResourceGroupName                                  string
	ApplicationGatewayWebApplicationFirewallPolicyName string
}

// NewApplicationGatewayWebApplicationFirewallPolicyID returns a new ApplicationGatewayWebApplicationFirewallPolicyId struct
func NewApplicationGatewayWebApplicationFirewallPolicyID(subscriptionId string, resourceGroupName string, applicationGatewayWebApplicationFirewallPolicyName string) ApplicationGatewayWebApplicationFirewallPolicyId {
	return ApplicationGatewayWebApplicationFirewallPolicyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ApplicationGatewayWebApplicationFirewallPolicyName: applicationGatewayWebApplicationFirewallPolicyName,
	}
}

// ParseApplicationGatewayWebApplicationFirewallPolicyID parses 'input' into a ApplicationGatewayWebApplicationFirewallPolicyId
func ParseApplicationGatewayWebApplicationFirewallPolicyID(input string) (*ApplicationGatewayWebApplicationFirewallPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationGatewayWebApplicationFirewallPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationGatewayWebApplicationFirewallPolicyId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively parses 'input' case-insensitively into a ApplicationGatewayWebApplicationFirewallPolicyId
// note: this method should only be used for API response data and not user input
func ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(input string) (*ApplicationGatewayWebApplicationFirewallPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationGatewayWebApplicationFirewallPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationGatewayWebApplicationFirewallPolicyId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApplicationGatewayWebApplicationFirewallPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ApplicationGatewayWebApplicationFirewallPolicyName, ok = input.Parsed["applicationGatewayWebApplicationFirewallPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "applicationGatewayWebApplicationFirewallPolicyName", input)
	}

	return nil
}

// ValidateApplicationGatewayWebApplicationFirewallPolicyID checks that 'input' can be parsed as a Application Gateway Web Application Firewall Policy ID
func ValidateApplicationGatewayWebApplicationFirewallPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationGatewayWebApplicationFirewallPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Gateway Web Application Firewall Policy ID
func (id ApplicationGatewayWebApplicationFirewallPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGatewayWebApplicationFirewallPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ApplicationGatewayWebApplicationFirewallPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Gateway Web Application Firewall Policy ID
func (id ApplicationGatewayWebApplicationFirewallPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticApplicationGatewayWebApplicationFirewallPolicies", "applicationGatewayWebApplicationFirewallPolicies", "applicationGatewayWebApplicationFirewallPolicies"),
		resourceids.UserSpecifiedSegment("applicationGatewayWebApplicationFirewallPolicyName", "applicationGatewayWebApplicationFirewallPolicyValue"),
	}
}

// String returns a human-readable description of this Application Gateway Web Application Firewall Policy ID
func (id ApplicationGatewayWebApplicationFirewallPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Application Gateway Web Application Firewall Policy Name: %q", id.ApplicationGatewayWebApplicationFirewallPolicyName),
	}
	return fmt.Sprintf("Application Gateway Web Application Firewall Policy (%s)", strings.Join(components, "\n"))
}