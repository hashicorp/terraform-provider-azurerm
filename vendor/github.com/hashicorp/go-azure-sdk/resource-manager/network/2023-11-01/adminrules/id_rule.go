package adminrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RuleId{})
}

var _ resourceids.ResourceId = &RuleId{}

// RuleId is a struct representing the Resource ID for a Rule
type RuleId struct {
	SubscriptionId                 string
	ResourceGroupName              string
	NetworkManagerName             string
	SecurityAdminConfigurationName string
	RuleCollectionName             string
	RuleName                       string
}

// NewRuleID returns a new RuleId struct
func NewRuleID(subscriptionId string, resourceGroupName string, networkManagerName string, securityAdminConfigurationName string, ruleCollectionName string, ruleName string) RuleId {
	return RuleId{
		SubscriptionId:                 subscriptionId,
		ResourceGroupName:              resourceGroupName,
		NetworkManagerName:             networkManagerName,
		SecurityAdminConfigurationName: securityAdminConfigurationName,
		RuleCollectionName:             ruleCollectionName,
		RuleName:                       ruleName,
	}
}

// ParseRuleID parses 'input' into a RuleId
func ParseRuleID(input string) (*RuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRuleIDInsensitively parses 'input' case-insensitively into a RuleId
// note: this method should only be used for API response data and not user input
func ParseRuleIDInsensitively(input string) (*RuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkManagerName, ok = input.Parsed["networkManagerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", input)
	}

	if id.SecurityAdminConfigurationName, ok = input.Parsed["securityAdminConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "securityAdminConfigurationName", input)
	}

	if id.RuleCollectionName, ok = input.Parsed["ruleCollectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "ruleCollectionName", input)
	}

	if id.RuleName, ok = input.Parsed["ruleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "ruleName", input)
	}

	return nil
}

// ValidateRuleID checks that 'input' can be parsed as a Rule ID
func ValidateRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Rule ID
func (id RuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/securityAdminConfigurations/%s/ruleCollections/%s/rules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName, id.RuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Rule ID
func (id RuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerName"),
		resourceids.StaticSegment("staticSecurityAdminConfigurations", "securityAdminConfigurations", "securityAdminConfigurations"),
		resourceids.UserSpecifiedSegment("securityAdminConfigurationName", "securityAdminConfigurationName"),
		resourceids.StaticSegment("staticRuleCollections", "ruleCollections", "ruleCollections"),
		resourceids.UserSpecifiedSegment("ruleCollectionName", "ruleCollectionName"),
		resourceids.StaticSegment("staticRules", "rules", "rules"),
		resourceids.UserSpecifiedSegment("ruleName", "ruleName"),
	}
}

// String returns a human-readable description of this Rule ID
func (id RuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
		fmt.Sprintf("Security Admin Configuration Name: %q", id.SecurityAdminConfigurationName),
		fmt.Sprintf("Rule Collection Name: %q", id.RuleCollectionName),
		fmt.Sprintf("Rule Name: %q", id.RuleName),
	}
	return fmt.Sprintf("Rule (%s)", strings.Join(components, "\n"))
}
