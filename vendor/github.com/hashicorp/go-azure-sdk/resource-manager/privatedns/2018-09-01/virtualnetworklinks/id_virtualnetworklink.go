package virtualnetworklinks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VirtualNetworkLinkId{}

// VirtualNetworkLinkId is a struct representing the Resource ID for a Virtual Network Link
type VirtualNetworkLinkId struct {
	SubscriptionId         string
	ResourceGroupName      string
	PrivateZoneName        string
	VirtualNetworkLinkName string
}

// NewVirtualNetworkLinkID returns a new VirtualNetworkLinkId struct
func NewVirtualNetworkLinkID(subscriptionId string, resourceGroupName string, privateZoneName string, virtualNetworkLinkName string) VirtualNetworkLinkId {
	return VirtualNetworkLinkId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		PrivateZoneName:        privateZoneName,
		VirtualNetworkLinkName: virtualNetworkLinkName,
	}
}

// ParseVirtualNetworkLinkID parses 'input' into a VirtualNetworkLinkId
func ParseVirtualNetworkLinkID(input string) (*VirtualNetworkLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualNetworkLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualNetworkLinkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PrivateZoneName, ok = parsed.Parsed["privateZoneName"]; !ok {
		return nil, fmt.Errorf("the segment 'privateZoneName' was not found in the resource id %q", input)
	}

	if id.VirtualNetworkLinkName, ok = parsed.Parsed["virtualNetworkLinkName"]; !ok {
		return nil, fmt.Errorf("the segment 'virtualNetworkLinkName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseVirtualNetworkLinkIDInsensitively parses 'input' case-insensitively into a VirtualNetworkLinkId
// note: this method should only be used for API response data and not user input
func ParseVirtualNetworkLinkIDInsensitively(input string) (*VirtualNetworkLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualNetworkLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualNetworkLinkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PrivateZoneName, ok = parsed.Parsed["privateZoneName"]; !ok {
		return nil, fmt.Errorf("the segment 'privateZoneName' was not found in the resource id %q", input)
	}

	if id.VirtualNetworkLinkName, ok = parsed.Parsed["virtualNetworkLinkName"]; !ok {
		return nil, fmt.Errorf("the segment 'virtualNetworkLinkName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateVirtualNetworkLinkID checks that 'input' can be parsed as a Virtual Network Link ID
func ValidateVirtualNetworkLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualNetworkLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Network Link ID
func (id VirtualNetworkLinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s/virtualNetworkLinks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateZoneName, id.VirtualNetworkLinkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Network Link ID
func (id VirtualNetworkLinkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticPrivateDnsZones", "privateDnsZones", "privateDnsZones"),
		resourceids.UserSpecifiedSegment("privateZoneName", "privateZoneValue"),
		resourceids.StaticSegment("staticVirtualNetworkLinks", "virtualNetworkLinks", "virtualNetworkLinks"),
		resourceids.UserSpecifiedSegment("virtualNetworkLinkName", "virtualNetworkLinkValue"),
	}
}

// String returns a human-readable description of this Virtual Network Link ID
func (id VirtualNetworkLinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Zone Name: %q", id.PrivateZoneName),
		fmt.Sprintf("Virtual Network Link Name: %q", id.VirtualNetworkLinkName),
	}
	return fmt.Sprintf("Virtual Network Link (%s)", strings.Join(components, "\n"))
}
