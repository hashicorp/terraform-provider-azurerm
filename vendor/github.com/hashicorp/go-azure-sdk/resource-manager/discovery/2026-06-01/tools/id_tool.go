package tools

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ToolId{})
}

var _ resourceids.ResourceId = &ToolId{}

// ToolId is a struct representing the Resource ID for a Tool
type ToolId struct {
	SubscriptionId    string
	ResourceGroupName string
	ToolName          string
}

// NewToolID returns a new ToolId struct
func NewToolID(subscriptionId string, resourceGroupName string, toolName string) ToolId {
	return ToolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ToolName:          toolName,
	}
}

// ParseToolID parses 'input' into a ToolId
func ParseToolID(input string) (*ToolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ToolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ToolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseToolIDInsensitively parses 'input' case-insensitively into a ToolId
// note: this method should only be used for API response data and not user input
func ParseToolIDInsensitively(input string) (*ToolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ToolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ToolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ToolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ToolName, ok = input.Parsed["toolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "toolName", input)
	}

	return nil
}

// ValidateToolID checks that 'input' can be parsed as a Tool ID
func ValidateToolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseToolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Tool ID
func (id ToolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Discovery/tools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ToolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Tool ID
func (id ToolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDiscovery", "Microsoft.Discovery", "Microsoft.Discovery"),
		resourceids.StaticSegment("staticTools", "tools", "tools"),
		resourceids.UserSpecifiedSegment("toolName", "toolName"),
	}
}

// String returns a human-readable description of this Tool ID
func (id ToolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Tool Name: %q", id.ToolName),
	}
	return fmt.Sprintf("Tool (%s)", strings.Join(components, "\n"))
}
