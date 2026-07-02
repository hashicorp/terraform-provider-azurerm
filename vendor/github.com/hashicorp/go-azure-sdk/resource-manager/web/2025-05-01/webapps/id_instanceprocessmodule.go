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
	recaser.RegisterResourceId(&InstanceProcessModuleId{})
}

var _ resourceids.ResourceId = &InstanceProcessModuleId{}

// InstanceProcessModuleId is a struct representing the Resource ID for a Instance Process Module
type InstanceProcessModuleId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	InstanceId        string
	ProcessId         string
	ModuleName        string
}

// NewInstanceProcessModuleID returns a new InstanceProcessModuleId struct
func NewInstanceProcessModuleID(subscriptionId string, resourceGroupName string, siteName string, instanceId string, processId string, moduleName string) InstanceProcessModuleId {
	return InstanceProcessModuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		InstanceId:        instanceId,
		ProcessId:         processId,
		ModuleName:        moduleName,
	}
}

// ParseInstanceProcessModuleID parses 'input' into a InstanceProcessModuleId
func ParseInstanceProcessModuleID(input string) (*InstanceProcessModuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InstanceProcessModuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InstanceProcessModuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseInstanceProcessModuleIDInsensitively parses 'input' case-insensitively into a InstanceProcessModuleId
// note: this method should only be used for API response data and not user input
func ParseInstanceProcessModuleIDInsensitively(input string) (*InstanceProcessModuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InstanceProcessModuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InstanceProcessModuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *InstanceProcessModuleId) FromParseResult(input resourceids.ParseResult) error {
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

// ValidateInstanceProcessModuleID checks that 'input' can be parsed as a Instance Process Module ID
func ValidateInstanceProcessModuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInstanceProcessModuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Instance Process Module ID
func (id InstanceProcessModuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/instances/%s/processes/%s/modules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.InstanceId, id.ProcessId, id.ModuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Instance Process Module ID
func (id InstanceProcessModuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticInstances", "instances", "instances"),
		resourceids.UserSpecifiedSegment("instanceId", "instanceId"),
		resourceids.StaticSegment("staticProcesses", "processes", "processes"),
		resourceids.UserSpecifiedSegment("processId", "processId"),
		resourceids.StaticSegment("staticModules", "modules", "modules"),
		resourceids.UserSpecifiedSegment("moduleName", "moduleName"),
	}
}

// String returns a human-readable description of this Instance Process Module ID
func (id InstanceProcessModuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Instance: %q", id.InstanceId),
		fmt.Sprintf("Process: %q", id.ProcessId),
		fmt.Sprintf("Module Name: %q", id.ModuleName),
	}
	return fmt.Sprintf("Instance Process Module (%s)", strings.Join(components, "\n"))
}
