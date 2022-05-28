package recordsets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = PrivateZoneId{}

// PrivateZoneId is a struct representing the Resource ID for a Private Zone
type PrivateZoneId struct {
	SubscriptionId    string
	ResourceGroupName string
	PrivateZoneName   string
	RecordType        RecordType
}

// NewPrivateZoneID returns a new PrivateZoneId struct
func NewPrivateZoneID(subscriptionId string, resourceGroupName string, privateZoneName string, recordType RecordType) PrivateZoneId {
	return PrivateZoneId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		PrivateZoneName:   privateZoneName,
		RecordType:        recordType,
	}
}

// ParsePrivateZoneID parses 'input' into a PrivateZoneId
func ParsePrivateZoneID(input string) (*PrivateZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateZoneId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateZoneId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PrivateZoneName, ok = parsed.Parsed["privateZoneName"]; !ok {
		return nil, fmt.Errorf("the segment 'privateZoneName' was not found in the resource id %q", input)
	}

	if v, constFound := parsed.Parsed["recordType"]; true {
		if !constFound {
			return nil, fmt.Errorf("the segment 'recordType' was not found in the resource id %q", input)
		}

		recordType, err := parseRecordType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.RecordType = *recordType
	}

	return &id, nil
}

// ParsePrivateZoneIDInsensitively parses 'input' case-insensitively into a PrivateZoneId
// note: this method should only be used for API response data and not user input
func ParsePrivateZoneIDInsensitively(input string) (*PrivateZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateZoneId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateZoneId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PrivateZoneName, ok = parsed.Parsed["privateZoneName"]; !ok {
		return nil, fmt.Errorf("the segment 'privateZoneName' was not found in the resource id %q", input)
	}

	if v, constFound := parsed.Parsed["recordType"]; true {
		if !constFound {
			return nil, fmt.Errorf("the segment 'recordType' was not found in the resource id %q", input)
		}

		recordType, err := parseRecordType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.RecordType = *recordType
	}

	return &id, nil
}

// ValidatePrivateZoneID checks that 'input' can be parsed as a Private Zone ID
func ValidatePrivateZoneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateZoneID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Zone ID
func (id PrivateZoneId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateZoneName, string(id.RecordType))
}

// Segments returns a slice of Resource ID Segments which comprise this Private Zone ID
func (id PrivateZoneId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticPrivateDnsZones", "privateDnsZones", "privateDnsZones"),
		resourceids.UserSpecifiedSegment("privateZoneName", "privateZoneValue"),
		resourceids.ConstantSegment("recordType", PossibleValuesForRecordType(), "A"),
	}
}

// String returns a human-readable description of this Private Zone ID
func (id PrivateZoneId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Zone Name: %q", id.PrivateZoneName),
		fmt.Sprintf("Record Type: %q", string(id.RecordType)),
	}
	return fmt.Sprintf("Private Zone (%s)", strings.Join(components, "\n"))
}
