package virtualharddisks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualHardDiskId{})
}

var _ resourceids.ResourceId = &VirtualHardDiskId{}

// VirtualHardDiskId is a struct representing the Resource ID for a Virtual Hard Disk
type VirtualHardDiskId struct {
	SubscriptionId      string
	ResourceGroupName   string
	VirtualHardDiskName string
}

// NewVirtualHardDiskID returns a new VirtualHardDiskId struct
func NewVirtualHardDiskID(subscriptionId string, resourceGroupName string, virtualHardDiskName string) VirtualHardDiskId {
	return VirtualHardDiskId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		VirtualHardDiskName: virtualHardDiskName,
	}
}

// ParseVirtualHardDiskID parses 'input' into a VirtualHardDiskId
func ParseVirtualHardDiskID(input string) (*VirtualHardDiskId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualHardDiskId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualHardDiskId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualHardDiskIDInsensitively parses 'input' case-insensitively into a VirtualHardDiskId
// note: this method should only be used for API response data and not user input
func ParseVirtualHardDiskIDInsensitively(input string) (*VirtualHardDiskId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualHardDiskId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualHardDiskId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualHardDiskId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualHardDiskName, ok = input.Parsed["virtualHardDiskName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualHardDiskName", input)
	}

	return nil
}

// ValidateVirtualHardDiskID checks that 'input' can be parsed as a Virtual Hard Disk ID
func ValidateVirtualHardDiskID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualHardDiskID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Hard Disk ID
func (id VirtualHardDiskId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/virtualHardDisks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualHardDiskName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Hard Disk ID
func (id VirtualHardDiskId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticVirtualHardDisks", "virtualHardDisks", "virtualHardDisks"),
		resourceids.UserSpecifiedSegment("virtualHardDiskName", "virtualHardDiskName"),
	}
}

// String returns a human-readable description of this Virtual Hard Disk ID
func (id VirtualHardDiskId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Hard Disk Name: %q", id.VirtualHardDiskName),
	}
	return fmt.Sprintf("Virtual Hard Disk (%s)", strings.Join(components, "\n"))
}
