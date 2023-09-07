package availabilitygrouplisteners

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AvailabilityGroupListenerId{}

// AvailabilityGroupListenerId is a struct representing the Resource ID for a Availability Group Listener
type AvailabilityGroupListenerId struct {
	SubscriptionId                string
	ResourceGroupName             string
	SqlVirtualMachineGroupName    string
	AvailabilityGroupListenerName string
}

// NewAvailabilityGroupListenerID returns a new AvailabilityGroupListenerId struct
func NewAvailabilityGroupListenerID(subscriptionId string, resourceGroupName string, sqlVirtualMachineGroupName string, availabilityGroupListenerName string) AvailabilityGroupListenerId {
	return AvailabilityGroupListenerId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		SqlVirtualMachineGroupName:    sqlVirtualMachineGroupName,
		AvailabilityGroupListenerName: availabilityGroupListenerName,
	}
}

// ParseAvailabilityGroupListenerID parses 'input' into a AvailabilityGroupListenerId
func ParseAvailabilityGroupListenerID(input string) (*AvailabilityGroupListenerId, error) {
	parser := resourceids.NewParserFromResourceIdType(AvailabilityGroupListenerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AvailabilityGroupListenerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SqlVirtualMachineGroupName, ok = parsed.Parsed["sqlVirtualMachineGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sqlVirtualMachineGroupName", *parsed)
	}

	if id.AvailabilityGroupListenerName, ok = parsed.Parsed["availabilityGroupListenerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "availabilityGroupListenerName", *parsed)
	}

	return &id, nil
}

// ParseAvailabilityGroupListenerIDInsensitively parses 'input' case-insensitively into a AvailabilityGroupListenerId
// note: this method should only be used for API response data and not user input
func ParseAvailabilityGroupListenerIDInsensitively(input string) (*AvailabilityGroupListenerId, error) {
	parser := resourceids.NewParserFromResourceIdType(AvailabilityGroupListenerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AvailabilityGroupListenerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SqlVirtualMachineGroupName, ok = parsed.Parsed["sqlVirtualMachineGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sqlVirtualMachineGroupName", *parsed)
	}

	if id.AvailabilityGroupListenerName, ok = parsed.Parsed["availabilityGroupListenerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "availabilityGroupListenerName", *parsed)
	}

	return &id, nil
}

// ValidateAvailabilityGroupListenerID checks that 'input' can be parsed as a Availability Group Listener ID
func ValidateAvailabilityGroupListenerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAvailabilityGroupListenerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Availability Group Listener ID
func (id AvailabilityGroupListenerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/%s/availabilityGroupListeners/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SqlVirtualMachineGroupName, id.AvailabilityGroupListenerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Availability Group Listener ID
func (id AvailabilityGroupListenerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSqlVirtualMachine", "Microsoft.SqlVirtualMachine", "Microsoft.SqlVirtualMachine"),
		resourceids.StaticSegment("staticSqlVirtualMachineGroups", "sqlVirtualMachineGroups", "sqlVirtualMachineGroups"),
		resourceids.UserSpecifiedSegment("sqlVirtualMachineGroupName", "sqlVirtualMachineGroupValue"),
		resourceids.StaticSegment("staticAvailabilityGroupListeners", "availabilityGroupListeners", "availabilityGroupListeners"),
		resourceids.UserSpecifiedSegment("availabilityGroupListenerName", "availabilityGroupListenerValue"),
	}
}

// String returns a human-readable description of this Availability Group Listener ID
func (id AvailabilityGroupListenerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Sql Virtual Machine Group Name: %q", id.SqlVirtualMachineGroupName),
		fmt.Sprintf("Availability Group Listener Name: %q", id.AvailabilityGroupListenerName),
	}
	return fmt.Sprintf("Availability Group Listener (%s)", strings.Join(components, "\n"))
}
