package virtualendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualEndpointId{})
}

var _ resourceids.ResourceId = &VirtualEndpointId{}

// VirtualEndpointId is a struct representing the Resource ID for a Virtual Endpoint
type VirtualEndpointId struct {
	SubscriptionId      string
	ResourceGroupName   string
	FlexibleServerName  string
	VirtualEndpointName string
}

// NewVirtualEndpointID returns a new VirtualEndpointId struct
func NewVirtualEndpointID(subscriptionId string, resourceGroupName string, flexibleServerName string, virtualEndpointName string) VirtualEndpointId {
	return VirtualEndpointId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		FlexibleServerName:  flexibleServerName,
		VirtualEndpointName: virtualEndpointName,
	}
}

// ParseVirtualEndpointID parses 'input' into a VirtualEndpointId
func ParseVirtualEndpointID(input string) (*VirtualEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualEndpointIDInsensitively parses 'input' case-insensitively into a VirtualEndpointId
// note: this method should only be used for API response data and not user input
func ParseVirtualEndpointIDInsensitively(input string) (*VirtualEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualEndpointId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FlexibleServerName, ok = input.Parsed["flexibleServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "flexibleServerName", input)
	}

	if id.VirtualEndpointName, ok = input.Parsed["virtualEndpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualEndpointName", input)
	}

	return nil
}

// ValidateVirtualEndpointID checks that 'input' can be parsed as a Virtual Endpoint ID
func ValidateVirtualEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Endpoint ID
func (id VirtualEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/flexibleServers/%s/virtualEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName, id.VirtualEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Endpoint ID
func (id VirtualEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforPostgreSQL", "Microsoft.DBforPostgreSQL", "Microsoft.DBforPostgreSQL"),
		resourceids.StaticSegment("staticFlexibleServers", "flexibleServers", "flexibleServers"),
		resourceids.UserSpecifiedSegment("flexibleServerName", "flexibleServerName"),
		resourceids.StaticSegment("staticVirtualEndpoints", "virtualEndpoints", "virtualEndpoints"),
		resourceids.UserSpecifiedSegment("virtualEndpointName", "virtualEndpointName"),
	}
}

// String returns a human-readable description of this Virtual Endpoint ID
func (id VirtualEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Flexible Server Name: %q", id.FlexibleServerName),
		fmt.Sprintf("Virtual Endpoint Name: %q", id.VirtualEndpointName),
	}
	return fmt.Sprintf("Virtual Endpoint (%s)", strings.Join(components, "\n"))
}
