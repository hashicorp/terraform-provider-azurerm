package adminrulecollections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RuleCollectionId{}

// RuleCollectionId is a struct representing the Resource ID for a Rule Collection
type RuleCollectionId struct {
	SubscriptionId                 string
	ResourceGroupName              string
	NetworkManagerName             string
	SecurityAdminConfigurationName string
	RuleCollectionName             string
}

// NewRuleCollectionID returns a new RuleCollectionId struct
func NewRuleCollectionID(subscriptionId string, resourceGroupName string, networkManagerName string, securityAdminConfigurationName string, ruleCollectionName string) RuleCollectionId {
	return RuleCollectionId{
		SubscriptionId:                 subscriptionId,
		ResourceGroupName:              resourceGroupName,
		NetworkManagerName:             networkManagerName,
		SecurityAdminConfigurationName: securityAdminConfigurationName,
		RuleCollectionName:             ruleCollectionName,
	}
}

// ParseRuleCollectionID parses 'input' into a RuleCollectionId
func ParseRuleCollectionID(input string) (*RuleCollectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(RuleCollectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RuleCollectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkManagerName, ok = parsed.Parsed["networkManagerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", *parsed)
	}

	if id.SecurityAdminConfigurationName, ok = parsed.Parsed["securityAdminConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "securityAdminConfigurationName", *parsed)
	}

	if id.RuleCollectionName, ok = parsed.Parsed["ruleCollectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ruleCollectionName", *parsed)
	}

	return &id, nil
}

// ParseRuleCollectionIDInsensitively parses 'input' case-insensitively into a RuleCollectionId
// note: this method should only be used for API response data and not user input
func ParseRuleCollectionIDInsensitively(input string) (*RuleCollectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(RuleCollectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RuleCollectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkManagerName, ok = parsed.Parsed["networkManagerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", *parsed)
	}

	if id.SecurityAdminConfigurationName, ok = parsed.Parsed["securityAdminConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "securityAdminConfigurationName", *parsed)
	}

	if id.RuleCollectionName, ok = parsed.Parsed["ruleCollectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ruleCollectionName", *parsed)
	}

	return &id, nil
}

// ValidateRuleCollectionID checks that 'input' can be parsed as a Rule Collection ID
func ValidateRuleCollectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRuleCollectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Rule Collection ID
func (id RuleCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/securityAdminConfigurations/%s/ruleCollections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Rule Collection ID
func (id RuleCollectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerValue"),
		resourceids.StaticSegment("staticSecurityAdminConfigurations", "securityAdminConfigurations", "securityAdminConfigurations"),
		resourceids.UserSpecifiedSegment("securityAdminConfigurationName", "securityAdminConfigurationValue"),
		resourceids.StaticSegment("staticRuleCollections", "ruleCollections", "ruleCollections"),
		resourceids.UserSpecifiedSegment("ruleCollectionName", "ruleCollectionValue"),
	}
}

// String returns a human-readable description of this Rule Collection ID
func (id RuleCollectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
		fmt.Sprintf("Security Admin Configuration Name: %q", id.SecurityAdminConfigurationName),
		fmt.Sprintf("Rule Collection Name: %q", id.RuleCollectionName),
	}
	return fmt.Sprintf("Rule Collection (%s)", strings.Join(components, "\n"))
}
