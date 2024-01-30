// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &VMwareSiteMachineId{}

// VMwareSiteMachineId is a struct representing the Resource ID for a VMware Site Machine
type VMwareSiteMachineId struct {
	SubscriptionId    string
	ResourceGroupName string
	VMwareSiteName    string
	MachineName       string
}

// NewVMwareSiteMachineID returns a new VMwareSiteMachineId struct
func NewVMwareSiteMachineID(subscriptionId string, resourceGroupName string, vmwareSiteName string, machineName string) VMwareSiteMachineId {
	return VMwareSiteMachineId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VMwareSiteName:    vmwareSiteName,
		MachineName:       machineName,
	}
}

// ParseVMwareSiteMachineID parses 'input' into a VMwareSiteMachineId
func ParseVMwareSiteMachineID(input string) (*VMwareSiteMachineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VMwareSiteMachineId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VMwareSiteMachineId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVMwareSiteMachineIDInsensitively parses 'input' case-insensitively into a VMwareSiteMachineId
// note: this method should only be used for API response data and not user input
func ParseVMwareSiteMachineIDInsensitively(input string) (*VMwareSiteMachineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VMwareSiteMachineId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VMwareSiteMachineId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VMwareSiteMachineId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VMwareSiteName, ok = input.Parsed["vmwareSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vmwareSiteName", input)
	}

	if id.MachineName, ok = input.Parsed["machineName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "machineName", input)
	}

	return nil
}

// ValidateVMwareSiteMachineID checks that 'input' can be parsed as a VMware Site Machine ID
func ValidateVMwareSiteMachineID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVMwareSiteMachineID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted VMware Site Machine ID
func (id VMwareSiteMachineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OffAzure/vmwareSites/%s/machines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VMwareSiteName, id.MachineName)
}

// Segments returns a slice of Resource ID Segments which comprise this VMware Site Machine ID
func (id VMwareSiteMachineId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOffAzure", "Microsoft.OffAzure", "Microsoft.OffAzure"),
		resourceids.StaticSegment("staticVMwareSites", "vmwareSites", "vmwareSites"),
		resourceids.UserSpecifiedSegment("vmwareSiteName", "vmwareSiteNameValue"),
		resourceids.StaticSegment("staticMachines", "machines", "machines"),
		resourceids.UserSpecifiedSegment("machineName", "machineValue"),
	}
}

// String returns a human-readable description of this VMware Site Machine ID
func (id VMwareSiteMachineId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("VMware Site Name: %q", id.VMwareSiteName),
		fmt.Sprintf("Machine Name: %q", id.MachineName),
	}
	return fmt.Sprintf("VMware Site Machine (%s)", strings.Join(components, "\n"))
}
