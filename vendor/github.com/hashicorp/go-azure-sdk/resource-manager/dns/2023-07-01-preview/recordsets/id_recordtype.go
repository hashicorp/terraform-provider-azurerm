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
	recaser.RegisterResourceId(&RecordTypeId{})
}

var _ resourceids.ResourceId = &RecordTypeId{}

// RecordTypeId is a struct representing the Resource ID for a Record Type
type RecordTypeId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DnsZoneName           string
	RecordType            RecordType
	RelativeRecordSetName string
}

// NewRecordTypeID returns a new RecordTypeId struct
func NewRecordTypeID(subscriptionId string, resourceGroupName string, dnsZoneName string, recordType RecordType, relativeRecordSetName string) RecordTypeId {
	return RecordTypeId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		DnsZoneName:           dnsZoneName,
		RecordType:            recordType,
		RelativeRecordSetName: relativeRecordSetName,
	}
}

// ParseRecordTypeID parses 'input' into a RecordTypeId
func ParseRecordTypeID(input string) (*RecordTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RecordTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RecordTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRecordTypeIDInsensitively parses 'input' case-insensitively into a RecordTypeId
// note: this method should only be used for API response data and not user input
func ParseRecordTypeIDInsensitively(input string) (*RecordTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RecordTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RecordTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RecordTypeId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.RelativeRecordSetName, ok = input.Parsed["relativeRecordSetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "relativeRecordSetName", input)
	}

	return nil
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsZones/%s/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName, string(id.RecordType), id.RelativeRecordSetName)
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
		resourceids.StaticSegment("staticDnsZones", "dnsZones", "dnsZones"),
		resourceids.UserSpecifiedSegment("dnsZoneName", "dnsZoneName"),
		resourceids.ConstantSegment("recordType", PossibleValuesForRecordType(), "A"),
		resourceids.UserSpecifiedSegment("relativeRecordSetName", "relativeRecordSetName"),
	}
}

// String returns a human-readable description of this Record Type ID
func (id RecordTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dns Zone Name: %q", id.DnsZoneName),
		fmt.Sprintf("Record Type: %q", string(id.RecordType)),
		fmt.Sprintf("Relative Record Set Name: %q", id.RelativeRecordSetName),
	}
	return fmt.Sprintf("Record Type (%s)", strings.Join(components, "\n"))
}
