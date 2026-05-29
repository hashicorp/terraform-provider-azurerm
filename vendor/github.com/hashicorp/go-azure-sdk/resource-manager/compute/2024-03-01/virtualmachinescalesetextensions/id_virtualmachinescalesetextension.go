package virtualmachinescalesetextensions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualMachineScaleSetExtensionId{})
}

var _ resourceids.ResourceId = &VirtualMachineScaleSetExtensionId{}

// VirtualMachineScaleSetExtensionId is a struct representing the Resource ID for a Virtual Machine Scale Set Extension
type VirtualMachineScaleSetExtensionId struct {
	SubscriptionId             string
	ResourceGroupName          string
	VirtualMachineScaleSetName string
	ExtensionName              string
}

// NewVirtualMachineScaleSetExtensionID returns a new VirtualMachineScaleSetExtensionId struct
func NewVirtualMachineScaleSetExtensionID(subscriptionId string, resourceGroupName string, virtualMachineScaleSetName string, extensionName string) VirtualMachineScaleSetExtensionId {
	return VirtualMachineScaleSetExtensionId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		VirtualMachineScaleSetName: virtualMachineScaleSetName,
		ExtensionName:              extensionName,
	}
}

// ParseVirtualMachineScaleSetExtensionID parses 'input' into a VirtualMachineScaleSetExtensionId
func ParseVirtualMachineScaleSetExtensionID(input string) (*VirtualMachineScaleSetExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineScaleSetExtensionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineScaleSetExtensionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualMachineScaleSetExtensionIDInsensitively parses 'input' case-insensitively into a VirtualMachineScaleSetExtensionId
// note: this method should only be used for API response data and not user input
func ParseVirtualMachineScaleSetExtensionIDInsensitively(input string) (*VirtualMachineScaleSetExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineScaleSetExtensionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineScaleSetExtensionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualMachineScaleSetExtensionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualMachineScaleSetName, ok = input.Parsed["virtualMachineScaleSetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualMachineScaleSetName", input)
	}

	if id.ExtensionName, ok = input.Parsed["extensionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "extensionName", input)
	}

	return nil
}

// ValidateVirtualMachineScaleSetExtensionID checks that 'input' can be parsed as a Virtual Machine Scale Set Extension ID
func ValidateVirtualMachineScaleSetExtensionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualMachineScaleSetExtensionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Machine Scale Set Extension ID
func (id VirtualMachineScaleSetExtensionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s/extensions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineScaleSetName, id.ExtensionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Machine Scale Set Extension ID
func (id VirtualMachineScaleSetExtensionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticVirtualMachineScaleSets", "virtualMachineScaleSets", "virtualMachineScaleSets"),
		resourceids.UserSpecifiedSegment("virtualMachineScaleSetName", "virtualMachineScaleSetName"),
		resourceids.StaticSegment("staticExtensions", "extensions", "extensions"),
		resourceids.UserSpecifiedSegment("extensionName", "extensionName"),
	}
}

// String returns a human-readable description of this Virtual Machine Scale Set Extension ID
func (id VirtualMachineScaleSetExtensionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Machine Scale Set Name: %q", id.VirtualMachineScaleSetName),
		fmt.Sprintf("Extension Name: %q", id.ExtensionName),
	}
	return fmt.Sprintf("Virtual Machine Scale Set Extension (%s)", strings.Join(components, "\n"))
}
