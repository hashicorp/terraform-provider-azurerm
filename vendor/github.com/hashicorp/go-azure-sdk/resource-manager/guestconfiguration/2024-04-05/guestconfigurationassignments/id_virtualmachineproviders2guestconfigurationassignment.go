package guestconfigurationassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualMachineProviders2GuestConfigurationAssignmentId{})
}

var _ resourceids.ResourceId = &VirtualMachineProviders2GuestConfigurationAssignmentId{}

// VirtualMachineProviders2GuestConfigurationAssignmentId is a struct representing the Resource ID for a Virtual Machine Providers 2 Guest Configuration Assignment
type VirtualMachineProviders2GuestConfigurationAssignmentId struct {
	SubscriptionId                   string
	ResourceGroupName                string
	VirtualMachineName               string
	GuestConfigurationAssignmentName string
}

// NewVirtualMachineProviders2GuestConfigurationAssignmentID returns a new VirtualMachineProviders2GuestConfigurationAssignmentId struct
func NewVirtualMachineProviders2GuestConfigurationAssignmentID(subscriptionId string, resourceGroupName string, virtualMachineName string, guestConfigurationAssignmentName string) VirtualMachineProviders2GuestConfigurationAssignmentId {
	return VirtualMachineProviders2GuestConfigurationAssignmentId{
		SubscriptionId:                   subscriptionId,
		ResourceGroupName:                resourceGroupName,
		VirtualMachineName:               virtualMachineName,
		GuestConfigurationAssignmentName: guestConfigurationAssignmentName,
	}
}

// ParseVirtualMachineProviders2GuestConfigurationAssignmentID parses 'input' into a VirtualMachineProviders2GuestConfigurationAssignmentId
func ParseVirtualMachineProviders2GuestConfigurationAssignmentID(input string) (*VirtualMachineProviders2GuestConfigurationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineProviders2GuestConfigurationAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineProviders2GuestConfigurationAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualMachineProviders2GuestConfigurationAssignmentIDInsensitively parses 'input' case-insensitively into a VirtualMachineProviders2GuestConfigurationAssignmentId
// note: this method should only be used for API response data and not user input
func ParseVirtualMachineProviders2GuestConfigurationAssignmentIDInsensitively(input string) (*VirtualMachineProviders2GuestConfigurationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineProviders2GuestConfigurationAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineProviders2GuestConfigurationAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualMachineProviders2GuestConfigurationAssignmentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualMachineName, ok = input.Parsed["virtualMachineName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualMachineName", input)
	}

	if id.GuestConfigurationAssignmentName, ok = input.Parsed["guestConfigurationAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "guestConfigurationAssignmentName", input)
	}

	return nil
}

// ValidateVirtualMachineProviders2GuestConfigurationAssignmentID checks that 'input' can be parsed as a Virtual Machine Providers 2 Guest Configuration Assignment ID
func ValidateVirtualMachineProviders2GuestConfigurationAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualMachineProviders2GuestConfigurationAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Machine Providers 2 Guest Configuration Assignment ID
func (id VirtualMachineProviders2GuestConfigurationAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineName, id.GuestConfigurationAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Machine Providers 2 Guest Configuration Assignment ID
func (id VirtualMachineProviders2GuestConfigurationAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticVirtualMachines", "virtualMachines", "virtualMachines"),
		resourceids.UserSpecifiedSegment("virtualMachineName", "virtualMachineName"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftGuestConfiguration", "Microsoft.GuestConfiguration", "Microsoft.GuestConfiguration"),
		resourceids.StaticSegment("staticGuestConfigurationAssignments", "guestConfigurationAssignments", "guestConfigurationAssignments"),
		resourceids.UserSpecifiedSegment("guestConfigurationAssignmentName", "guestConfigurationAssignmentName"),
	}
}

// String returns a human-readable description of this Virtual Machine Providers 2 Guest Configuration Assignment ID
func (id VirtualMachineProviders2GuestConfigurationAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Machine Name: %q", id.VirtualMachineName),
		fmt.Sprintf("Guest Configuration Assignment Name: %q", id.GuestConfigurationAssignmentName),
	}
	return fmt.Sprintf("Virtual Machine Providers 2 Guest Configuration Assignment (%s)", strings.Join(components, "\n"))
}
