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
	recaser.RegisterResourceId(&ProcessId{})
}

var _ resourceids.ResourceId = &ProcessId{}

// ProcessId is a struct representing the Resource ID for a Process
type ProcessId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	ProcessId         string
}

// NewProcessID returns a new ProcessId struct
func NewProcessID(subscriptionId string, resourceGroupName string, siteName string, processId string) ProcessId {
	return ProcessId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		ProcessId:         processId,
	}
}

// ParseProcessID parses 'input' into a ProcessId
func ParseProcessID(input string) (*ProcessId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProcessId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProcessId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProcessIDInsensitively parses 'input' case-insensitively into a ProcessId
// note: this method should only be used for API response data and not user input
func ParseProcessIDInsensitively(input string) (*ProcessId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProcessId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProcessId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProcessId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ProcessId, ok = input.Parsed["processId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "processId", input)
	}

	return nil
}

// ValidateProcessID checks that 'input' can be parsed as a Process ID
func ValidateProcessID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProcessID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Process ID
func (id ProcessId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/processes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.ProcessId)
}

// Segments returns a slice of Resource ID Segments which comprise this Process ID
func (id ProcessId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticProcesses", "processes", "processes"),
		resourceids.UserSpecifiedSegment("processId", "processId"),
	}
}

// String returns a human-readable description of this Process ID
func (id ProcessId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Process: %q", id.ProcessId),
	}
	return fmt.Sprintf("Process (%s)", strings.Join(components, "\n"))
}
