package virtualmachineruncommands

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualMachineRunCommandId{})
}

var _ resourceids.ResourceId = &VirtualMachineRunCommandId{}

// VirtualMachineRunCommandId is a struct representing the Resource ID for a Virtual Machine Run Command
type VirtualMachineRunCommandId struct {
	SubscriptionId     string
	ResourceGroupName  string
	VirtualMachineName string
	RunCommandName     string
}

// NewVirtualMachineRunCommandID returns a new VirtualMachineRunCommandId struct
func NewVirtualMachineRunCommandID(subscriptionId string, resourceGroupName string, virtualMachineName string, runCommandName string) VirtualMachineRunCommandId {
	return VirtualMachineRunCommandId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		VirtualMachineName: virtualMachineName,
		RunCommandName:     runCommandName,
	}
}

// ParseVirtualMachineRunCommandID parses 'input' into a VirtualMachineRunCommandId
func ParseVirtualMachineRunCommandID(input string) (*VirtualMachineRunCommandId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineRunCommandId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineRunCommandId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualMachineRunCommandIDInsensitively parses 'input' case-insensitively into a VirtualMachineRunCommandId
// note: this method should only be used for API response data and not user input
func ParseVirtualMachineRunCommandIDInsensitively(input string) (*VirtualMachineRunCommandId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineRunCommandId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineRunCommandId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualMachineRunCommandId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.RunCommandName, ok = input.Parsed["runCommandName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "runCommandName", input)
	}

	return nil
}

// ValidateVirtualMachineRunCommandID checks that 'input' can be parsed as a Virtual Machine Run Command ID
func ValidateVirtualMachineRunCommandID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualMachineRunCommandID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Machine Run Command ID
func (id VirtualMachineRunCommandId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/runCommands/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineName, id.RunCommandName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Machine Run Command ID
func (id VirtualMachineRunCommandId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticVirtualMachines", "virtualMachines", "virtualMachines"),
		resourceids.UserSpecifiedSegment("virtualMachineName", "virtualMachineName"),
		resourceids.StaticSegment("staticRunCommands", "runCommands", "runCommands"),
		resourceids.UserSpecifiedSegment("runCommandName", "runCommandName"),
	}
}

// String returns a human-readable description of this Virtual Machine Run Command ID
func (id VirtualMachineRunCommandId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Machine Name: %q", id.VirtualMachineName),
		fmt.Sprintf("Run Command Name: %q", id.RunCommandName),
	}
	return fmt.Sprintf("Virtual Machine Run Command (%s)", strings.Join(components, "\n"))
}
