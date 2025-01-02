package recordsets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ZoneId{})
}

var _ resourceids.ResourceId = &ZoneId{}

// ZoneId is a struct representing the Resource ID for a Zone
type ZoneId struct {
	SubscriptionId    string
	ResourceGroupName string
	DnsZoneName       string
	RecordType        RecordType
}

// NewZoneID returns a new ZoneId struct
func NewZoneID(subscriptionId string, resourceGroupName string, dnsZoneName string, recordType RecordType) ZoneId {
	return ZoneId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DnsZoneName:       dnsZoneName,
		RecordType:        recordType,
	}
}

// ParseZoneID parses 'input' into a ZoneId
func ParseZoneID(input string) (*ZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ZoneId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseZoneIDInsensitively parses 'input' case-insensitively into a ZoneId
// note: this method should only be used for API response data and not user input
func ParseZoneIDInsensitively(input string) (*ZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ZoneId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ZoneId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DnsZoneName, ok = input.Parsed["dnsZoneName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dnsZoneName", input)
	}

	if v, ok := input.Parsed["recordType"]; true {
		if !ok {
			return resourceids.NewSegmentNotSpecifiedError(id, "recordType", input)
		}

		recordType, err := parseRecordType(v)
		if err != nil {
			return fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.RecordType = *recordType
	}

	return nil
}

// ValidateZoneID checks that 'input' can be parsed as a Zone ID
func ValidateZoneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseZoneID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Zone ID
func (id ZoneId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsZones/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, string(id.RecordType))
}

// Segments returns a slice of Resource ID Segments which comprise this Zone ID
func (id ZoneId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDnsZones", "dnsZones", "dnsZones"),
		resourceids.UserSpecifiedSegment("dnsZoneName", "dnsZoneName"),
		resourceids.ConstantSegment("recordType", PossibleValuesForRecordType(), "A"),
	}
}

// String returns a human-readable description of this Zone ID
func (id ZoneId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dns Zone Name: %q", id.DnsZoneName),
		fmt.Sprintf("Record Type: %q", string(id.RecordType)),
	}
	return fmt.Sprintf("Zone (%s)", strings.Join(components, "\n"))
}
