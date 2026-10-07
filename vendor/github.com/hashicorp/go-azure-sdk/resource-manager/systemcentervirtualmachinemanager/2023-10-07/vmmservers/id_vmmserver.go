package vmmservers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VMmServerId{})
}

var _ resourceids.ResourceId = &VMmServerId{}

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
	parser := resourceids.NewParserFromResourceIdType(&VMmServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VMmServerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVMmServerIDInsensitively parses 'input' case-insensitively into a VMmServerId
// note: this method should only be used for API response data and not user input
func ParseVMmServerIDInsensitively(input string) (*VMmServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VMmServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VMmServerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VMmServerId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VmmServerName, ok = input.Parsed["vmmServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vmmServerName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("vmmServerName", "vmmServerName"),
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
