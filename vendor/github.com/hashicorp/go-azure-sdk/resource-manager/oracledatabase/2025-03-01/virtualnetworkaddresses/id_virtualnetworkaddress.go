package virtualnetworkaddresses

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualNetworkAddressId{})
}

var _ resourceids.ResourceId = &VirtualNetworkAddressId{}

// VirtualNetworkAddressId is a struct representing the Resource ID for a Virtual Network Address
type VirtualNetworkAddressId struct {
	SubscriptionId            string
	ResourceGroupName         string
	CloudVmClusterName        string
	VirtualNetworkAddressName string
}

// NewVirtualNetworkAddressID returns a new VirtualNetworkAddressId struct
func NewVirtualNetworkAddressID(subscriptionId string, resourceGroupName string, cloudVmClusterName string, virtualNetworkAddressName string) VirtualNetworkAddressId {
	return VirtualNetworkAddressId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		CloudVmClusterName:        cloudVmClusterName,
		VirtualNetworkAddressName: virtualNetworkAddressName,
	}
}

// ParseVirtualNetworkAddressID parses 'input' into a VirtualNetworkAddressId
func ParseVirtualNetworkAddressID(input string) (*VirtualNetworkAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkAddressId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkAddressId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualNetworkAddressIDInsensitively parses 'input' case-insensitively into a VirtualNetworkAddressId
// note: this method should only be used for API response data and not user input
func ParseVirtualNetworkAddressIDInsensitively(input string) (*VirtualNetworkAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkAddressId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkAddressId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualNetworkAddressId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CloudVmClusterName, ok = input.Parsed["cloudVmClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cloudVmClusterName", input)
	}

	if id.VirtualNetworkAddressName, ok = input.Parsed["virtualNetworkAddressName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkAddressName", input)
	}

	return nil
}

// ValidateVirtualNetworkAddressID checks that 'input' can be parsed as a Virtual Network Address ID
func ValidateVirtualNetworkAddressID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualNetworkAddressID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Network Address ID
func (id VirtualNetworkAddressId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/cloudVmClusters/%s/virtualNetworkAddresses/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudVmClusterName, id.VirtualNetworkAddressName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Network Address ID
func (id VirtualNetworkAddressId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticCloudVmClusters", "cloudVmClusters", "cloudVmClusters"),
		resourceids.UserSpecifiedSegment("cloudVmClusterName", "cloudVmClusterName"),
		resourceids.StaticSegment("staticVirtualNetworkAddresses", "virtualNetworkAddresses", "virtualNetworkAddresses"),
		resourceids.UserSpecifiedSegment("virtualNetworkAddressName", "virtualNetworkAddressName"),
	}
}

// String returns a human-readable description of this Virtual Network Address ID
func (id VirtualNetworkAddressId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Vm Cluster Name: %q", id.CloudVmClusterName),
		fmt.Sprintf("Virtual Network Address Name: %q", id.VirtualNetworkAddressName),
	}
	return fmt.Sprintf("Virtual Network Address (%s)", strings.Join(components, "\n"))
}
