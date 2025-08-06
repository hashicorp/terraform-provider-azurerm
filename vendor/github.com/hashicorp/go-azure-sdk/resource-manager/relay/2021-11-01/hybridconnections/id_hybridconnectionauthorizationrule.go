package hybridconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&HybridConnectionAuthorizationRuleId{})
}

var _ resourceids.ResourceId = &HybridConnectionAuthorizationRuleId{}

// HybridConnectionAuthorizationRuleId is a struct representing the Resource ID for a Hybrid Connection Authorization Rule
type HybridConnectionAuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NamespaceName         string
	HybridConnectionName  string
	AuthorizationRuleName string
}

// NewHybridConnectionAuthorizationRuleID returns a new HybridConnectionAuthorizationRuleId struct
func NewHybridConnectionAuthorizationRuleID(subscriptionId string, resourceGroupName string, namespaceName string, hybridConnectionName string, authorizationRuleName string) HybridConnectionAuthorizationRuleId {
	return HybridConnectionAuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NamespaceName:         namespaceName,
		HybridConnectionName:  hybridConnectionName,
		AuthorizationRuleName: authorizationRuleName,
	}
}

// ParseHybridConnectionAuthorizationRuleID parses 'input' into a HybridConnectionAuthorizationRuleId
func ParseHybridConnectionAuthorizationRuleID(input string) (*HybridConnectionAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HybridConnectionAuthorizationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HybridConnectionAuthorizationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHybridConnectionAuthorizationRuleIDInsensitively parses 'input' case-insensitively into a HybridConnectionAuthorizationRuleId
// note: this method should only be used for API response data and not user input
func ParseHybridConnectionAuthorizationRuleIDInsensitively(input string) (*HybridConnectionAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HybridConnectionAuthorizationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HybridConnectionAuthorizationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HybridConnectionAuthorizationRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NamespaceName, ok = input.Parsed["namespaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", input)
	}

	if id.HybridConnectionName, ok = input.Parsed["hybridConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hybridConnectionName", input)
	}

	if id.AuthorizationRuleName, ok = input.Parsed["authorizationRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authorizationRuleName", input)
	}

	return nil
}

// ValidateHybridConnectionAuthorizationRuleID checks that 'input' can be parsed as a Hybrid Connection Authorization Rule ID
func ValidateHybridConnectionAuthorizationRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHybridConnectionAuthorizationRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hybrid Connection Authorization Rule ID
func (id HybridConnectionAuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Relay/namespaces/%s/hybridConnections/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.HybridConnectionName, id.AuthorizationRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hybrid Connection Authorization Rule ID
func (id HybridConnectionAuthorizationRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRelay", "Microsoft.Relay", "Microsoft.Relay"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticHybridConnections", "hybridConnections", "hybridConnections"),
		resourceids.UserSpecifiedSegment("hybridConnectionName", "hybridConnectionName"),
		resourceids.StaticSegment("staticAuthorizationRules", "authorizationRules", "authorizationRules"),
		resourceids.UserSpecifiedSegment("authorizationRuleName", "authorizationRuleName"),
	}
}

// String returns a human-readable description of this Hybrid Connection Authorization Rule ID
func (id HybridConnectionAuthorizationRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Hybrid Connection Name: %q", id.HybridConnectionName),
		fmt.Sprintf("Authorization Rule Name: %q", id.AuthorizationRuleName),
	}
	return fmt.Sprintf("Hybrid Connection Authorization Rule (%s)", strings.Join(components, "\n"))
}
