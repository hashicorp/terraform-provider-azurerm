package inventoryitems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VMmServerId{}

// VMmServerId is a struct representing the Resource ID for a V Mm Server
type VMmServerId struct {
	SubscriptionId    string
	ResourceGroupName string
	VmmServerName     string
}

// NewVMmServerID returns a new VMmServerId struct
func NewVMmServerID(subscriptionId string, resourceGroupName string, vmmServerName string) VMmServerId {
	return VMmServerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VmmServerName:     vmmServerName,
	}
}

// ParseVMmServerID parses 'input' into a VMmServerId
func ParseVMmServerID(input string) (*VMmServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(VMmServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VMmServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VmmServerName, ok = parsed.Parsed["vmmServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vmmServerName", *parsed)
	}

	return &id, nil
}

// ParseVMmServerIDInsensitively parses 'input' case-insensitively into a VMmServerId
// note: this method should only be used for API response data and not user input
func ParseVMmServerIDInsensitively(input string) (*VMmServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(VMmServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VMmServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VmmServerName, ok = parsed.Parsed["vmmServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vmmServerName", *parsed)
	}

	return &id, nil
}

// ValidateVMmServerID checks that 'input' can be parsed as a V Mm Server ID
func ValidateVMmServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVMmServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted V Mm Server ID
func (id VMmServerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ScVmm/vmmServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VmmServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this V Mm Server ID
func (id VMmServerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftScVmm", "Microsoft.ScVmm", "Microsoft.ScVmm"),
		resourceids.StaticSegment("staticVmmServers", "vmmServers", "vmmServers"),
		resourceids.UserSpecifiedSegment("vmmServerName", "vmmServerValue"),
	}
}

// String returns a human-readable description of this V Mm Server ID
func (id VMmServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vmm Server Name: %q", id.VmmServerName),
	}
	return fmt.Sprintf("V Mm Server (%s)", strings.Join(components, "\n"))
}
