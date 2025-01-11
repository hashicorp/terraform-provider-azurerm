package configurationprofileassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualMachineProviders2ConfigurationProfileAssignmentId{})
}

var _ resourceids.ResourceId = &VirtualMachineProviders2ConfigurationProfileAssignmentId{}

// VirtualMachineProviders2ConfigurationProfileAssignmentId is a struct representing the Resource ID for a Virtual Machine Providers 2 Configuration Profile Assignment
type VirtualMachineProviders2ConfigurationProfileAssignmentId struct {
	SubscriptionId                     string
	ResourceGroupName                  string
	VirtualMachineName                 string
	ConfigurationProfileAssignmentName string
}

// NewVirtualMachineProviders2ConfigurationProfileAssignmentID returns a new VirtualMachineProviders2ConfigurationProfileAssignmentId struct
func NewVirtualMachineProviders2ConfigurationProfileAssignmentID(subscriptionId string, resourceGroupName string, virtualMachineName string, configurationProfileAssignmentName string) VirtualMachineProviders2ConfigurationProfileAssignmentId {
	return VirtualMachineProviders2ConfigurationProfileAssignmentId{
		SubscriptionId:                     subscriptionId,
		ResourceGroupName:                  resourceGroupName,
		VirtualMachineName:                 virtualMachineName,
		ConfigurationProfileAssignmentName: configurationProfileAssignmentName,
	}
}

// ParseVirtualMachineProviders2ConfigurationProfileAssignmentID parses 'input' into a VirtualMachineProviders2ConfigurationProfileAssignmentId
func ParseVirtualMachineProviders2ConfigurationProfileAssignmentID(input string) (*VirtualMachineProviders2ConfigurationProfileAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineProviders2ConfigurationProfileAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineProviders2ConfigurationProfileAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualMachineProviders2ConfigurationProfileAssignmentIDInsensitively parses 'input' case-insensitively into a VirtualMachineProviders2ConfigurationProfileAssignmentId
// note: this method should only be used for API response data and not user input
func ParseVirtualMachineProviders2ConfigurationProfileAssignmentIDInsensitively(input string) (*VirtualMachineProviders2ConfigurationProfileAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualMachineProviders2ConfigurationProfileAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualMachineProviders2ConfigurationProfileAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualMachineProviders2ConfigurationProfileAssignmentId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ConfigurationProfileAssignmentName, ok = input.Parsed["configurationProfileAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "configurationProfileAssignmentName", input)
	}

	return nil
}

// ValidateVirtualMachineProviders2ConfigurationProfileAssignmentID checks that 'input' can be parsed as a Virtual Machine Providers 2 Configuration Profile Assignment ID
func ValidateVirtualMachineProviders2ConfigurationProfileAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualMachineProviders2ConfigurationProfileAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Machine Providers 2 Configuration Profile Assignment ID
func (id VirtualMachineProviders2ConfigurationProfileAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/providers/Microsoft.AutoManage/configurationProfileAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineName, id.ConfigurationProfileAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Machine Providers 2 Configuration Profile Assignment ID
func (id VirtualMachineProviders2ConfigurationProfileAssignmentId) Segments() []resourceids.Segment {
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
		resourceids.ResourceProviderSegment("staticMicrosoftAutoManage", "Microsoft.AutoManage", "Microsoft.AutoManage"),
		resourceids.StaticSegment("staticConfigurationProfileAssignments", "configurationProfileAssignments", "configurationProfileAssignments"),
		resourceids.UserSpecifiedSegment("configurationProfileAssignmentName", "configurationProfileAssignmentName"),
	}
}

// String returns a human-readable description of this Virtual Machine Providers 2 Configuration Profile Assignment ID
func (id VirtualMachineProviders2ConfigurationProfileAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Machine Name: %q", id.VirtualMachineName),
		fmt.Sprintf("Configuration Profile Assignment Name: %q", id.ConfigurationProfileAssignmentName),
	}
	return fmt.Sprintf("Virtual Machine Providers 2 Configuration Profile Assignment (%s)", strings.Join(components, "\n"))
}
