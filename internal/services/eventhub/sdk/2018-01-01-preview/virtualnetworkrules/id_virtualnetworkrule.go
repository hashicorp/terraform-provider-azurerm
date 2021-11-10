package virtualnetworkrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VirtualnetworkruleId{}

// VirtualnetworkruleId is a struct representing the Resource ID for a Virtualnetworkrule
type VirtualnetworkruleId struct {
	SubscriptionId         string
	ResourceGroupName      string
	NamespaceName          string
	VirtualNetworkRuleName string
}

// NewVirtualnetworkruleID returns a new VirtualnetworkruleId struct
func NewVirtualnetworkruleID(subscriptionId string, resourceGroupName string, namespaceName string, virtualNetworkRuleName string) VirtualnetworkruleId {
	return VirtualnetworkruleId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		NamespaceName:          namespaceName,
		VirtualNetworkRuleName: virtualNetworkRuleName,
	}
}

// ParseVirtualnetworkruleID parses 'input' into a VirtualnetworkruleId
func ParseVirtualnetworkruleID(input string) (*VirtualnetworkruleId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualnetworkruleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualnetworkruleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.VirtualNetworkRuleName, ok = parsed.Parsed["virtualNetworkRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'virtualNetworkRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseVirtualnetworkruleIDInsensitively parses 'input' case-insensitively into a VirtualnetworkruleId
// note: this method should only be used for API response data and not user input
func ParseVirtualnetworkruleIDInsensitively(input string) (*VirtualnetworkruleId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualnetworkruleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualnetworkruleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.VirtualNetworkRuleName, ok = parsed.Parsed["virtualNetworkRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'virtualNetworkRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateVirtualnetworkruleID checks that 'input' can be parsed as a Virtualnetworkrule ID
func ValidateVirtualnetworkruleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualnetworkruleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtualnetworkrule ID
func (id VirtualnetworkruleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/virtualnetworkrules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.VirtualNetworkRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtualnetworkrule ID
func (id VirtualnetworkruleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("microsoftEventHub", "Microsoft.EventHub", "Microsoft.EventHub"),
		resourceids.StaticSegment("namespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("virtualnetworkrules", "virtualnetworkrules", "virtualnetworkrules"),
		resourceids.UserSpecifiedSegment("virtualNetworkRuleName", "virtualNetworkRuleValue"),
	}
}

// String returns a human-readable description of this Virtualnetworkrule ID
func (id VirtualnetworkruleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Virtual Network Rule Name: %q", id.VirtualNetworkRuleName),
	}
	return fmt.Sprintf("Virtualnetworkrule (%s)", strings.Join(components, "\n"))
}
