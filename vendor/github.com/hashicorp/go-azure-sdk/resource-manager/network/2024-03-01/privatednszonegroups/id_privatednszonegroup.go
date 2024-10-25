package privatednszonegroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PrivateDnsZoneGroupId{})
}

var _ resourceids.ResourceId = &PrivateDnsZoneGroupId{}

// PrivateDnsZoneGroupId is a struct representing the Resource ID for a Private Dns Zone Group
type PrivateDnsZoneGroupId struct {
	SubscriptionId          string
	ResourceGroupName       string
	PrivateEndpointName     string
	PrivateDnsZoneGroupName string
}

// NewPrivateDnsZoneGroupID returns a new PrivateDnsZoneGroupId struct
func NewPrivateDnsZoneGroupID(subscriptionId string, resourceGroupName string, privateEndpointName string, privateDnsZoneGroupName string) PrivateDnsZoneGroupId {
	return PrivateDnsZoneGroupId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		PrivateEndpointName:     privateEndpointName,
		PrivateDnsZoneGroupName: privateDnsZoneGroupName,
	}
}

// ParsePrivateDnsZoneGroupID parses 'input' into a PrivateDnsZoneGroupId
func ParsePrivateDnsZoneGroupID(input string) (*PrivateDnsZoneGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateDnsZoneGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateDnsZoneGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePrivateDnsZoneGroupIDInsensitively parses 'input' case-insensitively into a PrivateDnsZoneGroupId
// note: this method should only be used for API response data and not user input
func ParsePrivateDnsZoneGroupIDInsensitively(input string) (*PrivateDnsZoneGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateDnsZoneGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateDnsZoneGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PrivateDnsZoneGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PrivateEndpointName, ok = input.Parsed["privateEndpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointName", input)
	}

	if id.PrivateDnsZoneGroupName, ok = input.Parsed["privateDnsZoneGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateDnsZoneGroupName", input)
	}

	return nil
}

// ValidatePrivateDnsZoneGroupID checks that 'input' can be parsed as a Private Dns Zone Group ID
func ValidatePrivateDnsZoneGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateDnsZoneGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Dns Zone Group ID
func (id PrivateDnsZoneGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateEndpoints/%s/privateDnsZoneGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateEndpointName, id.PrivateDnsZoneGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Dns Zone Group ID
func (id PrivateDnsZoneGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticPrivateEndpoints", "privateEndpoints", "privateEndpoints"),
		resourceids.UserSpecifiedSegment("privateEndpointName", "privateEndpointName"),
		resourceids.StaticSegment("staticPrivateDnsZoneGroups", "privateDnsZoneGroups", "privateDnsZoneGroups"),
		resourceids.UserSpecifiedSegment("privateDnsZoneGroupName", "privateDnsZoneGroupName"),
	}
}

// String returns a human-readable description of this Private Dns Zone Group ID
func (id PrivateDnsZoneGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Endpoint Name: %q", id.PrivateEndpointName),
		fmt.Sprintf("Private Dns Zone Group Name: %q", id.PrivateDnsZoneGroupName),
	}
	return fmt.Sprintf("Private Dns Zone Group (%s)", strings.Join(components, "\n"))
}
