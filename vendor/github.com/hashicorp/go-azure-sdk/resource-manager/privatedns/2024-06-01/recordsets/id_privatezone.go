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
	recaser.RegisterResourceId(&PrivateZoneId{})
}

var _ resourceids.ResourceId = &PrivateZoneId{}

// PrivateZoneId is a struct representing the Resource ID for a Private Zone
type PrivateZoneId struct {
	SubscriptionId     string
	ResourceGroupName  string
	PrivateDnsZoneName string
	RecordType         RecordType
}

// NewPrivateZoneID returns a new PrivateZoneId struct
func NewPrivateZoneID(subscriptionId string, resourceGroupName string, privateDnsZoneName string, recordType RecordType) PrivateZoneId {
	return PrivateZoneId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		PrivateDnsZoneName: privateDnsZoneName,
		RecordType:         recordType,
	}
}

// ParsePrivateZoneID parses 'input' into a PrivateZoneId
func ParsePrivateZoneID(input string) (*PrivateZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateZoneId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePrivateZoneIDInsensitively parses 'input' case-insensitively into a PrivateZoneId
// note: this method should only be used for API response data and not user input
func ParsePrivateZoneIDInsensitively(input string) (*PrivateZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateZoneId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PrivateZoneId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PrivateDnsZoneName, ok = input.Parsed["privateDnsZoneName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateDnsZoneName", input)
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
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName, string(id.RecordType))
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
		resourceids.UserSpecifiedSegment("privateDnsZoneName", "privateDnsZoneName"),
		resourceids.ConstantSegment("recordType", PossibleValuesForRecordType(), "A"),
	}
}

// String returns a human-readable description of this Private Zone ID
func (id PrivateZoneId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Dns Zone Name: %q", id.PrivateDnsZoneName),
		fmt.Sprintf("Record Type: %q", string(id.RecordType)),
	}
	return fmt.Sprintf("Private Zone (%s)", strings.Join(components, "\n"))
}
