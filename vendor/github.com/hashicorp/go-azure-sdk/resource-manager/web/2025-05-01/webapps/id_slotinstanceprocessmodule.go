package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SlotInstanceProcessModuleId{})
}

var _ resourceids.ResourceId = &SlotInstanceProcessModuleId{}

// SlotInstanceProcessModuleId is a struct representing the Resource ID for a Slot Instance Process Module
type SlotInstanceProcessModuleId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SlotName          string
	InstanceId        string
	ProcessId         string
	ModuleName        string
}

// NewSlotInstanceProcessModuleID returns a new SlotInstanceProcessModuleId struct
func NewSlotInstanceProcessModuleID(subscriptionId string, resourceGroupName string, siteName string, slotName string, instanceId string, processId string, moduleName string) SlotInstanceProcessModuleId {
	return SlotInstanceProcessModuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SlotName:          slotName,
		InstanceId:        instanceId,
		ProcessId:         processId,
		ModuleName:        moduleName,
	}
}

// ParseSlotInstanceProcessModuleID parses 'input' into a SlotInstanceProcessModuleId
func ParseSlotInstanceProcessModuleID(input string) (*SlotInstanceProcessModuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotInstanceProcessModuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotInstanceProcessModuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotInstanceProcessModuleIDInsensitively parses 'input' case-insensitively into a SlotInstanceProcessModuleId
// note: this method should only be used for API response data and not user input
func ParseSlotInstanceProcessModuleIDInsensitively(input string) (*SlotInstanceProcessModuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotInstanceProcessModuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotInstanceProcessModuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotInstanceProcessModuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.SlotName, ok = input.Parsed["slotName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "slotName", input)
	}

	if id.InstanceId, ok = input.Parsed["instanceId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "instanceId", input)
	}

	if id.ProcessId, ok = input.Parsed["processId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "processId", input)
	}

	if id.ModuleName, ok = input.Parsed["moduleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "moduleName", input)
	}

	return nil
}

// ValidateSlotInstanceProcessModuleID checks that 'input' can be parsed as a Slot Instance Process Module ID
func ValidateSlotInstanceProcessModuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotInstanceProcessModuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Instance Process Module ID
func (id SlotInstanceProcessModuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/instances/%s/processes/%s/modules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.InstanceId, id.ProcessId, id.ModuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Instance Process Module ID
func (id SlotInstanceProcessModuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticSlots", "slots", "slots"),
		resourceids.UserSpecifiedSegment("slotName", "slotName"),
		resourceids.StaticSegment("staticInstances", "instances", "instances"),
		resourceids.UserSpecifiedSegment("instanceId", "instanceId"),
		resourceids.StaticSegment("staticProcesses", "processes", "processes"),
		resourceids.UserSpecifiedSegment("processId", "processId"),
		resourceids.StaticSegment("staticModules", "modules", "modules"),
		resourceids.UserSpecifiedSegment("moduleName", "moduleName"),
	}
}

// String returns a human-readable description of this Slot Instance Process Module ID
func (id SlotInstanceProcessModuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Instance: %q", id.InstanceId),
		fmt.Sprintf("Process: %q", id.ProcessId),
		fmt.Sprintf("Module Name: %q", id.ModuleName),
	}
	return fmt.Sprintf("Slot Instance Process Module (%s)", strings.Join(components, "\n"))
}
