package connectiongateways

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ConnectionGatewayId{}

// ConnectionGatewayId is a struct representing the Resource ID for a Connection Gateway
type ConnectionGatewayId struct {
	SubscriptionId        string
	ResourceGroupName     string
	ConnectionGatewayName string
}

// NewConnectionGatewayID returns a new ConnectionGatewayId struct
func NewConnectionGatewayID(subscriptionId string, resourceGroupName string, connectionGatewayName string) ConnectionGatewayId {
	return ConnectionGatewayId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		ConnectionGatewayName: connectionGatewayName,
	}
}

// ParseConnectionGatewayID parses 'input' into a ConnectionGatewayId
func ParseConnectionGatewayID(input string) (*ConnectionGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConnectionGatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConnectionGatewayId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ConnectionGatewayName, ok = parsed.Parsed["connectionGatewayName"]; !ok {
		return nil, fmt.Errorf("the segment 'connectionGatewayName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseConnectionGatewayIDInsensitively parses 'input' case-insensitively into a ConnectionGatewayId
// note: this method should only be used for API response data and not user input
func ParseConnectionGatewayIDInsensitively(input string) (*ConnectionGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConnectionGatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConnectionGatewayId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ConnectionGatewayName, ok = parsed.Parsed["connectionGatewayName"]; !ok {
		return nil, fmt.Errorf("the segment 'connectionGatewayName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateConnectionGatewayID checks that 'input' can be parsed as a Connection Gateway ID
func ValidateConnectionGatewayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectionGatewayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connection Gateway ID
func (id ConnectionGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/connectionGateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ConnectionGatewayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connection Gateway ID
func (id ConnectionGatewayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticConnectionGateways", "connectionGateways", "connectionGateways"),
		resourceids.UserSpecifiedSegment("connectionGatewayName", "connectionGatewayValue"),
	}
}

// String returns a human-readable description of this Connection Gateway ID
func (id ConnectionGatewayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Connection Gateway Name: %q", id.ConnectionGatewayName),
	}
	return fmt.Sprintf("Connection Gateway (%s)", strings.Join(components, "\n"))
}
