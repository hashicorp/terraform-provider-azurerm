package hybridconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = HybridConnectionId{}

// HybridConnectionId is a struct representing the Resource ID for a Hybrid Connection
type HybridConnectionId struct {
	SubscriptionId       string
	ResourceGroupName    string
	NamespaceName        string
	HybridConnectionName string
}

// NewHybridConnectionID returns a new HybridConnectionId struct
func NewHybridConnectionID(subscriptionId string, resourceGroupName string, namespaceName string, hybridConnectionName string) HybridConnectionId {
	return HybridConnectionId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		NamespaceName:        namespaceName,
		HybridConnectionName: hybridConnectionName,
	}
}

// ParseHybridConnectionID parses 'input' into a HybridConnectionId
func ParseHybridConnectionID(input string) (*HybridConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(HybridConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HybridConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.HybridConnectionName, ok = parsed.Parsed["hybridConnectionName"]; !ok {
		return nil, fmt.Errorf("the segment 'hybridConnectionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseHybridConnectionIDInsensitively parses 'input' case-insensitively into a HybridConnectionId
// note: this method should only be used for API response data and not user input
func ParseHybridConnectionIDInsensitively(input string) (*HybridConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(HybridConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HybridConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.HybridConnectionName, ok = parsed.Parsed["hybridConnectionName"]; !ok {
		return nil, fmt.Errorf("the segment 'hybridConnectionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateHybridConnectionID checks that 'input' can be parsed as a Hybrid Connection ID
func ValidateHybridConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHybridConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hybrid Connection ID
func (id HybridConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Relay/namespaces/%s/hybridConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.HybridConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hybrid Connection ID
func (id HybridConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRelay", "Microsoft.Relay", "Microsoft.Relay"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticHybridConnections", "hybridConnections", "hybridConnections"),
		resourceids.UserSpecifiedSegment("hybridConnectionName", "hybridConnectionValue"),
	}
}

// String returns a human-readable description of this Hybrid Connection ID
func (id HybridConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Hybrid Connection Name: %q", id.HybridConnectionName),
	}
	return fmt.Sprintf("Hybrid Connection (%s)", strings.Join(components, "\n"))
}
