package recordsets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = RecordTypeId{}

// RecordTypeId is a struct representing the Resource ID for a Record Type
type RecordTypeId struct {
	SubscriptionId        string
	ResourceGroupName     string
	PrivateZoneName       string
	RecordType            RecordType
	RelativeRecordSetName string
}

// NewRecordTypeID returns a new RecordTypeId struct
func NewRecordTypeID(subscriptionId string, resourceGroupName string, privateZoneName string, recordType RecordType, relativeRecordSetName string) RecordTypeId {
	return RecordTypeId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		PrivateZoneName:       privateZoneName,
		RecordType:            recordType,
		RelativeRecordSetName: relativeRecordSetName,
	}
}

// ParseRecordTypeID parses 'input' into a RecordTypeId
func ParseRecordTypeID(input string) (*RecordTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(RecordTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RecordTypeId{}

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

	if id.RelativeRecordSetName, ok = parsed.Parsed["relativeRecordSetName"]; !ok {
		return nil, fmt.Errorf("the segment 'relativeRecordSetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseRecordTypeIDInsensitively parses 'input' case-insensitively into a RecordTypeId
// note: this method should only be used for API response data and not user input
func ParseRecordTypeIDInsensitively(input string) (*RecordTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(RecordTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RecordTypeId{}

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

	if id.RelativeRecordSetName, ok = parsed.Parsed["relativeRecordSetName"]; !ok {
		return nil, fmt.Errorf("the segment 'relativeRecordSetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateRecordTypeID checks that 'input' can be parsed as a Record Type ID
func ValidateRecordTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRecordTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Record Type ID
func (id RecordTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateZoneName, string(id.RecordType), id.RelativeRecordSetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Record Type ID
func (id RecordTypeId) Segments() []resourceids.Segment {
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
		resourceids.UserSpecifiedSegment("relativeRecordSetName", "relativeRecordSetValue"),
	}
}

// String returns a human-readable description of this Record Type ID
func (id RecordTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Zone Name: %q", id.PrivateZoneName),
		fmt.Sprintf("Record Type: %q", string(id.RecordType)),
		fmt.Sprintf("Relative Record Set Name: %q", id.RelativeRecordSetName),
	}
	return fmt.Sprintf("Record Type (%s)", strings.Join(components, "\n"))
}
