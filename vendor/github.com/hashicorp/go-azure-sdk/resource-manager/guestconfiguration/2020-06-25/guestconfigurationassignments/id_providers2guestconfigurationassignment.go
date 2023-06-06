package guestconfigurationassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = Providers2GuestConfigurationAssignmentId{}

// Providers2GuestConfigurationAssignmentId is a struct representing the Resource ID for a Providers 2 Guest Configuration Assignment
type Providers2GuestConfigurationAssignmentId struct {
	SubscriptionId                   string
	ResourceGroupName                string
	VirtualMachineName               string
	GuestConfigurationAssignmentName string
}

// NewProviders2GuestConfigurationAssignmentID returns a new Providers2GuestConfigurationAssignmentId struct
func NewProviders2GuestConfigurationAssignmentID(subscriptionId string, resourceGroupName string, virtualMachineName string, guestConfigurationAssignmentName string) Providers2GuestConfigurationAssignmentId {
	return Providers2GuestConfigurationAssignmentId{
		SubscriptionId:                   subscriptionId,
		ResourceGroupName:                resourceGroupName,
		VirtualMachineName:               virtualMachineName,
		GuestConfigurationAssignmentName: guestConfigurationAssignmentName,
	}
}

// ParseProviders2GuestConfigurationAssignmentID parses 'input' into a Providers2GuestConfigurationAssignmentId
func ParseProviders2GuestConfigurationAssignmentID(input string) (*Providers2GuestConfigurationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(Providers2GuestConfigurationAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := Providers2GuestConfigurationAssignmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualMachineName, ok = parsed.Parsed["virtualMachineName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualMachineName", *parsed)
	}

	if id.GuestConfigurationAssignmentName, ok = parsed.Parsed["guestConfigurationAssignmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "guestConfigurationAssignmentName", *parsed)
	}

	return &id, nil
}

// ParseProviders2GuestConfigurationAssignmentIDInsensitively parses 'input' case-insensitively into a Providers2GuestConfigurationAssignmentId
// note: this method should only be used for API response data and not user input
func ParseProviders2GuestConfigurationAssignmentIDInsensitively(input string) (*Providers2GuestConfigurationAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(Providers2GuestConfigurationAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := Providers2GuestConfigurationAssignmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualMachineName, ok = parsed.Parsed["virtualMachineName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualMachineName", *parsed)
	}

	if id.GuestConfigurationAssignmentName, ok = parsed.Parsed["guestConfigurationAssignmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "guestConfigurationAssignmentName", *parsed)
	}

	return &id, nil
}

// ValidateProviders2GuestConfigurationAssignmentID checks that 'input' can be parsed as a Providers 2 Guest Configuration Assignment ID
func ValidateProviders2GuestConfigurationAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviders2GuestConfigurationAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Providers 2 Guest Configuration Assignment ID
func (id Providers2GuestConfigurationAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineName, id.GuestConfigurationAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Providers 2 Guest Configuration Assignment ID
func (id Providers2GuestConfigurationAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticVirtualMachines", "virtualMachines", "virtualMachines"),
		resourceids.UserSpecifiedSegment("virtualMachineName", "virtualMachineValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftGuestConfiguration", "Microsoft.GuestConfiguration", "Microsoft.GuestConfiguration"),
		resourceids.StaticSegment("staticGuestConfigurationAssignments", "guestConfigurationAssignments", "guestConfigurationAssignments"),
		resourceids.UserSpecifiedSegment("guestConfigurationAssignmentName", "guestConfigurationAssignmentValue"),
	}
}

// String returns a human-readable description of this Providers 2 Guest Configuration Assignment ID
func (id Providers2GuestConfigurationAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Machine Name: %q", id.VirtualMachineName),
		fmt.Sprintf("Guest Configuration Assignment Name: %q", id.GuestConfigurationAssignmentName),
	}
	return fmt.Sprintf("Providers 2 Guest Configuration Assignment (%s)", strings.Join(components, "\n"))
}
