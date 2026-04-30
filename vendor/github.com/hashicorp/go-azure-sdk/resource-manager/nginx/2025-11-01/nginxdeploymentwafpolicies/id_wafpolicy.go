package nginxdeploymentwafpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WafPolicyId{})
}

var _ resourceids.ResourceId = &WafPolicyId{}

// WafPolicyId is a struct representing the Resource ID for a Waf Policy
type WafPolicyId struct {
	SubscriptionId      string
	ResourceGroupName   string
	NginxDeploymentName string
	WafPolicyName       string
}

// NewWafPolicyID returns a new WafPolicyId struct
func NewWafPolicyID(subscriptionId string, resourceGroupName string, nginxDeploymentName string, wafPolicyName string) WafPolicyId {
	return WafPolicyId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		NginxDeploymentName: nginxDeploymentName,
		WafPolicyName:       wafPolicyName,
	}
}

// ParseWafPolicyID parses 'input' into a WafPolicyId
func ParseWafPolicyID(input string) (*WafPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WafPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WafPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWafPolicyIDInsensitively parses 'input' case-insensitively into a WafPolicyId
// note: this method should only be used for API response data and not user input
func ParseWafPolicyIDInsensitively(input string) (*WafPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WafPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WafPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WafPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NginxDeploymentName, ok = input.Parsed["nginxDeploymentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "nginxDeploymentName", input)
	}

	if id.WafPolicyName, ok = input.Parsed["wafPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "wafPolicyName", input)
	}

	return nil
}

// ValidateWafPolicyID checks that 'input' can be parsed as a Waf Policy ID
func ValidateWafPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWafPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Waf Policy ID
func (id WafPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Nginx.NginxPlus/nginxDeployments/%s/wafPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName, id.WafPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Waf Policy ID
func (id WafPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticNginxNginxPlus", "Nginx.NginxPlus", "Nginx.NginxPlus"),
		resourceids.StaticSegment("staticNginxDeployments", "nginxDeployments", "nginxDeployments"),
		resourceids.UserSpecifiedSegment("nginxDeploymentName", "nginxDeploymentName"),
		resourceids.StaticSegment("staticWafPolicies", "wafPolicies", "wafPolicies"),
		resourceids.UserSpecifiedSegment("wafPolicyName", "wafPolicyName"),
	}
}

// String returns a human-readable description of this Waf Policy ID
func (id WafPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Nginx Deployment Name: %q", id.NginxDeploymentName),
		fmt.Sprintf("Waf Policy Name: %q", id.WafPolicyName),
	}
	return fmt.Sprintf("Waf Policy (%s)", strings.Join(components, "\n"))
}
