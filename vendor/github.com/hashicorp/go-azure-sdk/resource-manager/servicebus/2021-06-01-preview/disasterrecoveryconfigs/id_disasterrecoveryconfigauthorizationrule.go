package disasterrecoveryconfigs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DisasterRecoveryConfigAuthorizationRuleId{}

// DisasterRecoveryConfigAuthorizationRuleId is a struct representing the Resource ID for a Disaster Recovery Config Authorization Rule
type DisasterRecoveryConfigAuthorizationRuleId struct {
	SubscriptionId             string
	ResourceGroupName          string
	NamespaceName              string
	DisasterRecoveryConfigName string
	AuthorizationRuleName      string
}

// NewDisasterRecoveryConfigAuthorizationRuleID returns a new DisasterRecoveryConfigAuthorizationRuleId struct
func NewDisasterRecoveryConfigAuthorizationRuleID(subscriptionId string, resourceGroupName string, namespaceName string, disasterRecoveryConfigName string, authorizationRuleName string) DisasterRecoveryConfigAuthorizationRuleId {
	return DisasterRecoveryConfigAuthorizationRuleId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		NamespaceName:              namespaceName,
		DisasterRecoveryConfigName: disasterRecoveryConfigName,
		AuthorizationRuleName:      authorizationRuleName,
	}
}

// ParseDisasterRecoveryConfigAuthorizationRuleID parses 'input' into a DisasterRecoveryConfigAuthorizationRuleId
func ParseDisasterRecoveryConfigAuthorizationRuleID(input string) (*DisasterRecoveryConfigAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(DisasterRecoveryConfigAuthorizationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DisasterRecoveryConfigAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.DisasterRecoveryConfigName, ok = parsed.Parsed["disasterRecoveryConfigName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "disasterRecoveryConfigName", *parsed)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "authorizationRuleName", *parsed)
	}

	return &id, nil
}

// ParseDisasterRecoveryConfigAuthorizationRuleIDInsensitively parses 'input' case-insensitively into a DisasterRecoveryConfigAuthorizationRuleId
// note: this method should only be used for API response data and not user input
func ParseDisasterRecoveryConfigAuthorizationRuleIDInsensitively(input string) (*DisasterRecoveryConfigAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(DisasterRecoveryConfigAuthorizationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DisasterRecoveryConfigAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.DisasterRecoveryConfigName, ok = parsed.Parsed["disasterRecoveryConfigName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "disasterRecoveryConfigName", *parsed)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "authorizationRuleName", *parsed)
	}

	return &id, nil
}

// ValidateDisasterRecoveryConfigAuthorizationRuleID checks that 'input' can be parsed as a Disaster Recovery Config Authorization Rule ID
func ValidateDisasterRecoveryConfigAuthorizationRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDisasterRecoveryConfigAuthorizationRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Disaster Recovery Config Authorization Rule ID
func (id DisasterRecoveryConfigAuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/disasterRecoveryConfigs/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.DisasterRecoveryConfigName, id.AuthorizationRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Disaster Recovery Config Authorization Rule ID
func (id DisasterRecoveryConfigAuthorizationRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceBus", "Microsoft.ServiceBus", "Microsoft.ServiceBus"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticDisasterRecoveryConfigs", "disasterRecoveryConfigs", "disasterRecoveryConfigs"),
		resourceids.UserSpecifiedSegment("disasterRecoveryConfigName", "disasterRecoveryConfigValue"),
		resourceids.StaticSegment("staticAuthorizationRules", "authorizationRules", "authorizationRules"),
		resourceids.UserSpecifiedSegment("authorizationRuleName", "authorizationRuleValue"),
	}
}

// String returns a human-readable description of this Disaster Recovery Config Authorization Rule ID
func (id DisasterRecoveryConfigAuthorizationRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Disaster Recovery Config Name: %q", id.DisasterRecoveryConfigName),
		fmt.Sprintf("Authorization Rule Name: %q", id.AuthorizationRuleName),
	}
	return fmt.Sprintf("Disaster Recovery Config Authorization Rule (%s)", strings.Join(components, "\n"))
}
