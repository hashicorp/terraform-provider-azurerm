package wcfrelays

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = WcfRelayAuthorizationRuleId{}

// WcfRelayAuthorizationRuleId is a struct representing the Resource ID for a Wcf Relay Authorization Rule
type WcfRelayAuthorizationRuleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NamespaceName         string
	RelayName             string
	AuthorizationRuleName string
}

// NewWcfRelayAuthorizationRuleID returns a new WcfRelayAuthorizationRuleId struct
func NewWcfRelayAuthorizationRuleID(subscriptionId string, resourceGroupName string, namespaceName string, relayName string, authorizationRuleName string) WcfRelayAuthorizationRuleId {
	return WcfRelayAuthorizationRuleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NamespaceName:         namespaceName,
		RelayName:             relayName,
		AuthorizationRuleName: authorizationRuleName,
	}
}

// ParseWcfRelayAuthorizationRuleID parses 'input' into a WcfRelayAuthorizationRuleId
func ParseWcfRelayAuthorizationRuleID(input string) (*WcfRelayAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(WcfRelayAuthorizationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WcfRelayAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.RelayName, ok = parsed.Parsed["relayName"]; !ok {
		return nil, fmt.Errorf("the segment 'relayName' was not found in the resource id %q", input)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'authorizationRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseWcfRelayAuthorizationRuleIDInsensitively parses 'input' case-insensitively into a WcfRelayAuthorizationRuleId
// note: this method should only be used for API response data and not user input
func ParseWcfRelayAuthorizationRuleIDInsensitively(input string) (*WcfRelayAuthorizationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(WcfRelayAuthorizationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WcfRelayAuthorizationRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.RelayName, ok = parsed.Parsed["relayName"]; !ok {
		return nil, fmt.Errorf("the segment 'relayName' was not found in the resource id %q", input)
	}

	if id.AuthorizationRuleName, ok = parsed.Parsed["authorizationRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'authorizationRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateWcfRelayAuthorizationRuleID checks that 'input' can be parsed as a Wcf Relay Authorization Rule ID
func ValidateWcfRelayAuthorizationRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWcfRelayAuthorizationRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Wcf Relay Authorization Rule ID
func (id WcfRelayAuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Relay/namespaces/%s/wcfRelays/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.RelayName, id.AuthorizationRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Wcf Relay Authorization Rule ID
func (id WcfRelayAuthorizationRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("microsoftRelay", "Microsoft.Relay", "Microsoft.Relay"),
		resourceids.StaticSegment("namespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("wcfRelays", "wcfRelays", "wcfRelays"),
		resourceids.UserSpecifiedSegment("relayName", "relayValue"),
		resourceids.StaticSegment("authorizationRules", "authorizationRules", "authorizationRules"),
		resourceids.UserSpecifiedSegment("authorizationRuleName", "authorizationRuleValue"),
	}
}

// String returns a human-readable description of this Wcf Relay Authorization Rule ID
func (id WcfRelayAuthorizationRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Relay Name: %q", id.RelayName),
		fmt.Sprintf("Authorization Rule Name: %q", id.AuthorizationRuleName),
	}
	return fmt.Sprintf("Wcf Relay Authorization Rule (%s)", strings.Join(components, "\n"))
}
