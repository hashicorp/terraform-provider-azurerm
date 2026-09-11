package pipelinegroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PipelineGroupId{})
}

var _ resourceids.ResourceId = &PipelineGroupId{}

// PipelineGroupId is a struct representing the Resource ID for a Pipeline Group
type PipelineGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	PipelineGroupName string
}

// NewPipelineGroupID returns a new PipelineGroupId struct
func NewPipelineGroupID(subscriptionId string, resourceGroupName string, pipelineGroupName string) PipelineGroupId {
	return PipelineGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		PipelineGroupName: pipelineGroupName,
	}
}

// ParsePipelineGroupID parses 'input' into a PipelineGroupId
func ParsePipelineGroupID(input string) (*PipelineGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PipelineGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PipelineGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePipelineGroupIDInsensitively parses 'input' case-insensitively into a PipelineGroupId
// note: this method should only be used for API response data and not user input
func ParsePipelineGroupIDInsensitively(input string) (*PipelineGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PipelineGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PipelineGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PipelineGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PipelineGroupName, ok = input.Parsed["pipelineGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "pipelineGroupName", input)
	}

	return nil
}

// ValidatePipelineGroupID checks that 'input' can be parsed as a Pipeline Group ID
func ValidatePipelineGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePipelineGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Pipeline Group ID
func (id PipelineGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Monitor/pipelineGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PipelineGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Pipeline Group ID
func (id PipelineGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMonitor", "Microsoft.Monitor", "Microsoft.Monitor"),
		resourceids.StaticSegment("staticPipelineGroups", "pipelineGroups", "pipelineGroups"),
		resourceids.UserSpecifiedSegment("pipelineGroupName", "pipelineGroupName"),
	}
}

// String returns a human-readable description of this Pipeline Group ID
func (id PipelineGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Pipeline Group Name: %q", id.PipelineGroupName),
	}
	return fmt.Sprintf("Pipeline Group (%s)", strings.Join(components, "\n"))
}
