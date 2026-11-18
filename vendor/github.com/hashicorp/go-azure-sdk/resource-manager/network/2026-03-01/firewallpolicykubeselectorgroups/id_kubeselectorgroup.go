package firewallpolicykubeselectorgroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&KubeSelectorGroupId{})
}

var _ resourceids.ResourceId = &KubeSelectorGroupId{}

// KubeSelectorGroupId is a struct representing the Resource ID for a Kube Selector Group
type KubeSelectorGroupId struct {
	SubscriptionId        string
	ResourceGroupName     string
	FirewallPolicyName    string
	KubeSelectorGroupName string
}

// NewKubeSelectorGroupID returns a new KubeSelectorGroupId struct
func NewKubeSelectorGroupID(subscriptionId string, resourceGroupName string, firewallPolicyName string, kubeSelectorGroupName string) KubeSelectorGroupId {
	return KubeSelectorGroupId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		FirewallPolicyName:    firewallPolicyName,
		KubeSelectorGroupName: kubeSelectorGroupName,
	}
}

// ParseKubeSelectorGroupID parses 'input' into a KubeSelectorGroupId
func ParseKubeSelectorGroupID(input string) (*KubeSelectorGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KubeSelectorGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KubeSelectorGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseKubeSelectorGroupIDInsensitively parses 'input' case-insensitively into a KubeSelectorGroupId
// note: this method should only be used for API response data and not user input
func ParseKubeSelectorGroupIDInsensitively(input string) (*KubeSelectorGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KubeSelectorGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KubeSelectorGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *KubeSelectorGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FirewallPolicyName, ok = input.Parsed["firewallPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "firewallPolicyName", input)
	}

	if id.KubeSelectorGroupName, ok = input.Parsed["kubeSelectorGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "kubeSelectorGroupName", input)
	}

	return nil
}

// ValidateKubeSelectorGroupID checks that 'input' can be parsed as a Kube Selector Group ID
func ValidateKubeSelectorGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKubeSelectorGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Kube Selector Group ID
func (id KubeSelectorGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/firewallPolicies/%s/kubeSelectorGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FirewallPolicyName, id.KubeSelectorGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Kube Selector Group ID
func (id KubeSelectorGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticFirewallPolicies", "firewallPolicies", "firewallPolicies"),
		resourceids.UserSpecifiedSegment("firewallPolicyName", "firewallPolicyName"),
		resourceids.StaticSegment("staticKubeSelectorGroups", "kubeSelectorGroups", "kubeSelectorGroups"),
		resourceids.UserSpecifiedSegment("kubeSelectorGroupName", "kubeSelectorGroupName"),
	}
}

// String returns a human-readable description of this Kube Selector Group ID
func (id KubeSelectorGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Firewall Policy Name: %q", id.FirewallPolicyName),
		fmt.Sprintf("Kube Selector Group Name: %q", id.KubeSelectorGroupName),
	}
	return fmt.Sprintf("Kube Selector Group (%s)", strings.Join(components, "\n"))
}
