package pipelineruns

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PipelineRunId{})
}

var _ resourceids.ResourceId = &PipelineRunId{}

// PipelineRunId is a struct representing the Resource ID for a Pipeline Run
type PipelineRunId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	PipelineRunName   string
}

// NewPipelineRunID returns a new PipelineRunId struct
func NewPipelineRunID(subscriptionId string, resourceGroupName string, registryName string, pipelineRunName string) PipelineRunId {
	return PipelineRunId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		PipelineRunName:   pipelineRunName,
	}
}

// ParsePipelineRunID parses 'input' into a PipelineRunId
func ParsePipelineRunID(input string) (*PipelineRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PipelineRunId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PipelineRunId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePipelineRunIDInsensitively parses 'input' case-insensitively into a PipelineRunId
// note: this method should only be used for API response data and not user input
func ParsePipelineRunIDInsensitively(input string) (*PipelineRunId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PipelineRunId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PipelineRunId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PipelineRunId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RegistryName, ok = input.Parsed["registryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "registryName", input)
	}

	if id.PipelineRunName, ok = input.Parsed["pipelineRunName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "pipelineRunName", input)
	}

	return nil
}

// ValidatePipelineRunID checks that 'input' can be parsed as a Pipeline Run ID
func ValidatePipelineRunID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePipelineRunID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Pipeline Run ID
func (id PipelineRunId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/pipelineRuns/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.PipelineRunName)
}

// Segments returns a slice of Resource ID Segments which comprise this Pipeline Run ID
func (id PipelineRunId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryName"),
		resourceids.StaticSegment("staticPipelineRuns", "pipelineRuns", "pipelineRuns"),
		resourceids.UserSpecifiedSegment("pipelineRunName", "pipelineRunName"),
	}
}

// String returns a human-readable description of this Pipeline Run ID
func (id PipelineRunId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Pipeline Run Name: %q", id.PipelineRunName),
	}
	return fmt.Sprintf("Pipeline Run (%s)", strings.Join(components, "\n"))
}
