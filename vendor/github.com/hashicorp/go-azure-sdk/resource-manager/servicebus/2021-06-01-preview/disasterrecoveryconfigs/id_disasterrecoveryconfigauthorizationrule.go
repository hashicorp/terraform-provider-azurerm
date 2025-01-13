package disasterrecoveryconfigs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DisasterRecoveryConfigAuthorizationRuleId{})
}

var _ resourceids.ResourceId = &DisasterRecoveryConfigAuthorizationRuleId{}

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
	parser := resourceids.NewParserFromResourceIdType(&DisasterRecoveryConfigAuthorizationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DisasterRecoveryConfigAuthorizationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDisasterRecoveryConfigAuthorizationRuleIDInsensitively parses 'input' case-insensitively into a DisasterRecoveryConfigAuthorizationRuleId
// note: this method should only be used for API response data and not user input
func ParseDisasterRecoveryConfigAuthorizationRuleIDInsensitively(input string) (*DisasterRecoveryConfigAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DisasterRecoveryConfigAuthorizationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DisasterRecoveryConfigAuthorizationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DisasterRecoveryConfigAuthorizationRuleId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DisasterRecoveryConfigName, ok = input.Parsed["disasterRecoveryConfigName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "disasterRecoveryConfigName", input)
	}

	if id.AuthorizationRuleName, ok = input.Parsed["authorizationRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authorizationRuleName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticDisasterRecoveryConfigs", "disasterRecoveryConfigs", "disasterRecoveryConfigs"),
		resourceids.UserSpecifiedSegment("disasterRecoveryConfigName", "disasterRecoveryConfigName"),
		resourceids.StaticSegment("staticAuthorizationRules", "authorizationRules", "authorizationRules"),
		resourceids.UserSpecifiedSegment("authorizationRuleName", "authorizationRuleName"),
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
