package standbyvirtualmachinepools

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StandbyVirtualMachinePoolId{})
}

var _ resourceids.ResourceId = &StandbyVirtualMachinePoolId{}

// StandbyVirtualMachinePoolId is a struct representing the Resource ID for a Standby Virtual Machine Pool
type StandbyVirtualMachinePoolId struct {
	SubscriptionId                string
	ResourceGroupName             string
	StandbyVirtualMachinePoolName string
}

// NewStandbyVirtualMachinePoolID returns a new StandbyVirtualMachinePoolId struct
func NewStandbyVirtualMachinePoolID(subscriptionId string, resourceGroupName string, standbyVirtualMachinePoolName string) StandbyVirtualMachinePoolId {
	return StandbyVirtualMachinePoolId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		StandbyVirtualMachinePoolName: standbyVirtualMachinePoolName,
	}
}

// ParseStandbyVirtualMachinePoolID parses 'input' into a StandbyVirtualMachinePoolId
func ParseStandbyVirtualMachinePoolID(input string) (*StandbyVirtualMachinePoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StandbyVirtualMachinePoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StandbyVirtualMachinePoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStandbyVirtualMachinePoolIDInsensitively parses 'input' case-insensitively into a StandbyVirtualMachinePoolId
// note: this method should only be used for API response data and not user input
func ParseStandbyVirtualMachinePoolIDInsensitively(input string) (*StandbyVirtualMachinePoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StandbyVirtualMachinePoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StandbyVirtualMachinePoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StandbyVirtualMachinePoolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StandbyVirtualMachinePoolName, ok = input.Parsed["standbyVirtualMachinePoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "standbyVirtualMachinePoolName", input)
	}

	return nil
}

// ValidateStandbyVirtualMachinePoolID checks that 'input' can be parsed as a Standby Virtual Machine Pool ID
func ValidateStandbyVirtualMachinePoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStandbyVirtualMachinePoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Standby Virtual Machine Pool ID
func (id StandbyVirtualMachinePoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StandbyPool/standbyVirtualMachinePools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StandbyVirtualMachinePoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Standby Virtual Machine Pool ID
func (id StandbyVirtualMachinePoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStandbyPool", "Microsoft.StandbyPool", "Microsoft.StandbyPool"),
		resourceids.StaticSegment("staticStandbyVirtualMachinePools", "standbyVirtualMachinePools", "standbyVirtualMachinePools"),
		resourceids.UserSpecifiedSegment("standbyVirtualMachinePoolName", "standbyVirtualMachinePoolName"),
	}
}

// String returns a human-readable description of this Standby Virtual Machine Pool ID
func (id StandbyVirtualMachinePoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Standby Virtual Machine Pool Name: %q", id.StandbyVirtualMachinePoolName),
	}
	return fmt.Sprintf("Standby Virtual Machine Pool (%s)", strings.Join(components, "\n"))
}
