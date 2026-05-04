package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualNetworkConnectionId{})
}

var _ resourceids.ResourceId = &VirtualNetworkConnectionId{}

// VirtualNetworkConnectionId is a struct representing the Resource ID for a Virtual Network Connection
type VirtualNetworkConnectionId struct {
	SubscriptionId               string
	ResourceGroupName            string
	SiteName                     string
	VirtualNetworkConnectionName string
}

// NewVirtualNetworkConnectionID returns a new VirtualNetworkConnectionId struct
func NewVirtualNetworkConnectionID(subscriptionId string, resourceGroupName string, siteName string, virtualNetworkConnectionName string) VirtualNetworkConnectionId {
	return VirtualNetworkConnectionId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		SiteName:                     siteName,
		VirtualNetworkConnectionName: virtualNetworkConnectionName,
	}
}

// ParseVirtualNetworkConnectionID parses 'input' into a VirtualNetworkConnectionId
func ParseVirtualNetworkConnectionID(input string) (*VirtualNetworkConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualNetworkConnectionIDInsensitively parses 'input' case-insensitively into a VirtualNetworkConnectionId
// note: this method should only be used for API response data and not user input
func ParseVirtualNetworkConnectionIDInsensitively(input string) (*VirtualNetworkConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualNetworkConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.VirtualNetworkConnectionName, ok = input.Parsed["virtualNetworkConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkConnectionName", input)
	}

	return nil
}

// ValidateVirtualNetworkConnectionID checks that 'input' can be parsed as a Virtual Network Connection ID
func ValidateVirtualNetworkConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualNetworkConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Network Connection ID
func (id VirtualNetworkConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/virtualNetworkConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.VirtualNetworkConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Network Connection ID
func (id VirtualNetworkConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticVirtualNetworkConnections", "virtualNetworkConnections", "virtualNetworkConnections"),
		resourceids.UserSpecifiedSegment("virtualNetworkConnectionName", "virtualNetworkConnectionName"),
	}
}

// String returns a human-readable description of this Virtual Network Connection ID
func (id VirtualNetworkConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Virtual Network Connection Name: %q", id.VirtualNetworkConnectionName),
	}
	return fmt.Sprintf("Virtual Network Connection (%s)", strings.Join(components, "\n"))
}
