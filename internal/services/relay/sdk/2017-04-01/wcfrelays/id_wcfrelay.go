package wcfrelays

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = WcfRelayId{}

// WcfRelayId is a struct representing the Resource ID for a Wcf Relay
type WcfRelayId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	RelayName         string
}

// NewWcfRelayID returns a new WcfRelayId struct
func NewWcfRelayID(subscriptionId string, resourceGroupName string, namespaceName string, relayName string) WcfRelayId {
	return WcfRelayId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		RelayName:         relayName,
	}
}

// ParseWcfRelayID parses 'input' into a WcfRelayId
func ParseWcfRelayID(input string) (*WcfRelayId, error) {
	parser := resourceids.NewParserFromResourceIdType(WcfRelayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WcfRelayId{}

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

	return &id, nil
}

// ParseWcfRelayIDInsensitively parses 'input' case-insensitively into a WcfRelayId
// note: this method should only be used for API response data and not user input
func ParseWcfRelayIDInsensitively(input string) (*WcfRelayId, error) {
	parser := resourceids.NewParserFromResourceIdType(WcfRelayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WcfRelayId{}

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

	return &id, nil
}

// ValidateWcfRelayID checks that 'input' can be parsed as a Wcf Relay ID
func ValidateWcfRelayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWcfRelayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Wcf Relay ID
func (id WcfRelayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Relay/namespaces/%s/wcfRelays/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.RelayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Wcf Relay ID
func (id WcfRelayId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Wcf Relay ID
func (id WcfRelayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Relay Name: %q", id.RelayName),
	}
	return fmt.Sprintf("Wcf Relay (%s)", strings.Join(components, "\n"))
}
