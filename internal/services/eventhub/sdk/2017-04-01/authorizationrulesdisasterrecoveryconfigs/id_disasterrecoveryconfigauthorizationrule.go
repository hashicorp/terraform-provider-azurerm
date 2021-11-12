package authorizationrulesdisasterrecoveryconfigs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DisasterRecoveryConfigAuthorizationRuleId{}

// DisasterRecoveryConfigAuthorizationRuleId is a struct representing the Resource ID for a Disaster Recovery Config Authorization Rule
type DisasterRecoveryConfigAuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NamespaceName         string
	Alias                 string
	AuthorizationRuleName string
}

// NewDisasterRecoveryConfigAuthorizationRuleID returns a new DisasterRecoveryConfigAuthorizationRuleId struct
func NewDisasterRecoveryConfigAuthorizationRuleID(subscriptionId string, resourceGroupName string, namespaceName string, alias string, authorizationRuleName string) DisasterRecoveryConfigAuthorizationRuleId {
	return DisasterRecoveryConfigAuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NamespaceName:         namespaceName,
		Alias:                 alias,
		AuthorizationRuleName: authorizationRuleName,
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
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.Alias, ok = parsed.Parsed["alias"]; !ok {
		return nil, fmt.Errorf("the segment 'alias' was not found in the resource id %q", input)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'authorizationRuleName' was not found in the resource id %q", input)
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
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.Alias, ok = parsed.Parsed["alias"]; !ok {
		return nil, fmt.Errorf("the segment 'alias' was not found in the resource id %q", input)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'authorizationRuleName' was not found in the resource id %q", input)
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/disasterRecoveryConfigs/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.Alias, id.AuthorizationRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Disaster Recovery Config Authorization Rule ID
func (id DisasterRecoveryConfigAuthorizationRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventHub", "Microsoft.EventHub", "Microsoft.EventHub"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticDisasterRecoveryConfigs", "disasterRecoveryConfigs", "disasterRecoveryConfigs"),
		resourceids.UserSpecifiedSegment("alias", "aliasValue"),
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
		fmt.Sprintf("Alias: %q", id.Alias),
		fmt.Sprintf("Authorization Rule Name: %q", id.AuthorizationRuleName),
	}
	return fmt.Sprintf("Disaster Recovery Config Authorization Rule (%s)", strings.Join(components, "\n"))
}
