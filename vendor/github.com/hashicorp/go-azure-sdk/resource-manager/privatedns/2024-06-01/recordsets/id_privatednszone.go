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
	recaser.RegisterResourceId(&PrivateDnsZoneId{})
}

var _ resourceids.ResourceId = &PrivateDnsZoneId{}

// PrivateDnsZoneId is a struct representing the Resource ID for a Private Dns Zone
type PrivateDnsZoneId struct {
	SubscriptionId     string
	ResourceGroupName  string
	PrivateDnsZoneName string
}

// NewPrivateDnsZoneID returns a new PrivateDnsZoneId struct
func NewPrivateDnsZoneID(subscriptionId string, resourceGroupName string, privateDnsZoneName string) PrivateDnsZoneId {
	return PrivateDnsZoneId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		PrivateDnsZoneName: privateDnsZoneName,
	}
}

// ParsePrivateDnsZoneID parses 'input' into a PrivateDnsZoneId
func ParsePrivateDnsZoneID(input string) (*PrivateDnsZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateDnsZoneId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateDnsZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePrivateDnsZoneIDInsensitively parses 'input' case-insensitively into a PrivateDnsZoneId
// note: this method should only be used for API response data and not user input
func ParsePrivateDnsZoneIDInsensitively(input string) (*PrivateDnsZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateDnsZoneId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateDnsZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PrivateDnsZoneId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
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
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName)
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
		resourceids.UserSpecifiedSegment("privateDnsZoneName", "privateDnsZoneName"),
	}
}

// String returns a human-readable description of this Private Dns Zone ID
func (id PrivateDnsZoneId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Dns Zone Name: %q", id.PrivateDnsZoneName),
	}
	return fmt.Sprintf("Private Dns Zone (%s)", strings.Join(components, "\n"))
}
