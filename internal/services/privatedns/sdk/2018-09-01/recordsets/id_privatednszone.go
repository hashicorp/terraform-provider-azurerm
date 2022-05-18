package recordsets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = PrivateDnsZoneId{}

// PrivateDnsZoneId is a struct representing the Resource ID for a Private Dns Zone
type PrivateDnsZoneId struct {
	SubscriptionId    string
	ResourceGroupName string
	PrivateZoneName   string
}

// NewPrivateDnsZoneID returns a new PrivateDnsZoneId struct
func NewPrivateDnsZoneID(subscriptionId string, resourceGroupName string, privateZoneName string) PrivateDnsZoneId {
	return PrivateDnsZoneId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		PrivateZoneName:   privateZoneName,
	}
}

// ParsePrivateDnsZoneID parses 'input' into a PrivateDnsZoneId
func ParsePrivateDnsZoneID(input string) (*PrivateDnsZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateDnsZoneId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateDnsZoneId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PrivateZoneName, ok = parsed.Parsed["privateZoneName"]; !ok {
		return nil, fmt.Errorf("the segment 'privateZoneName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParsePrivateDnsZoneIDInsensitively parses 'input' case-insensitively into a PrivateDnsZoneId
// note: this method should only be used for API response data and not user input
func ParsePrivateDnsZoneIDInsensitively(input string) (*PrivateDnsZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateDnsZoneId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateDnsZoneId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PrivateZoneName, ok = parsed.Parsed["privateZoneName"]; !ok {
		return nil, fmt.Errorf("the segment 'privateZoneName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidatePrivateDnsZoneID checks that 'input' can be parsed as a Private Dns Zone ID
func ValidatePrivateDnsZoneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateDnsZoneID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Dns Zone ID
func (id PrivateDnsZoneId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateZoneName)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Dns Zone ID
func (id PrivateDnsZoneId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticPrivateDnsZones", "privateDnsZones", "privateDnsZones"),
		resourceids.UserSpecifiedSegment("privateZoneName", "privateZoneValue"),
	}
}

// String returns a human-readable description of this Private Dns Zone ID
func (id PrivateDnsZoneId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Zone Name: %q", id.PrivateZoneName),
	}
	return fmt.Sprintf("Private Dns Zone (%s)", strings.Join(components, "\n"))
}
