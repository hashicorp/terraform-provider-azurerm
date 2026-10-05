package managedprivateendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ManagedVirtualNetworkId{})
}

var _ resourceids.ResourceId = &ManagedVirtualNetworkId{}

// ManagedVirtualNetworkId is a struct representing the Resource ID for a Managed Virtual Network
type ManagedVirtualNetworkId struct {
	SubscriptionId            string
	ResourceGroupName         string
	FactoryName               string
	ManagedVirtualNetworkName string
}

// NewManagedVirtualNetworkID returns a new ManagedVirtualNetworkId struct
func NewManagedVirtualNetworkID(subscriptionId string, resourceGroupName string, factoryName string, managedVirtualNetworkName string) ManagedVirtualNetworkId {
	return ManagedVirtualNetworkId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		FactoryName:               factoryName,
		ManagedVirtualNetworkName: managedVirtualNetworkName,
	}
}

// ParseManagedVirtualNetworkID parses 'input' into a ManagedVirtualNetworkId
func ParseManagedVirtualNetworkID(input string) (*ManagedVirtualNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedVirtualNetworkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedVirtualNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseManagedVirtualNetworkIDInsensitively parses 'input' case-insensitively into a ManagedVirtualNetworkId
// note: this method should only be used for API response data and not user input
func ParseManagedVirtualNetworkIDInsensitively(input string) (*ManagedVirtualNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedVirtualNetworkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedVirtualNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ManagedVirtualNetworkId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FactoryName, ok = input.Parsed["factoryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "factoryName", input)
	}

	if id.ManagedVirtualNetworkName, ok = input.Parsed["managedVirtualNetworkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedVirtualNetworkName", input)
	}

	return nil
}

// ValidateManagedVirtualNetworkID checks that 'input' can be parsed as a Managed Virtual Network ID
func ValidateManagedVirtualNetworkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedVirtualNetworkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Virtual Network ID
func (id ManagedVirtualNetworkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s/managedVirtualNetworks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FactoryName, id.ManagedVirtualNetworkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Virtual Network ID
func (id ManagedVirtualNetworkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataFactory", "Microsoft.DataFactory", "Microsoft.DataFactory"),
		resourceids.StaticSegment("staticFactories", "factories", "factories"),
		resourceids.UserSpecifiedSegment("factoryName", "factoryName"),
		resourceids.StaticSegment("staticManagedVirtualNetworks", "managedVirtualNetworks", "managedVirtualNetworks"),
		resourceids.UserSpecifiedSegment("managedVirtualNetworkName", "managedVirtualNetworkName"),
	}
}

// String returns a human-readable description of this Managed Virtual Network ID
func (id ManagedVirtualNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Factory Name: %q", id.FactoryName),
		fmt.Sprintf("Managed Virtual Network Name: %q", id.ManagedVirtualNetworkName),
	}
	return fmt.Sprintf("Managed Virtual Network (%s)", strings.Join(components, "\n"))
}
