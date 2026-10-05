package runs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RunId{})
}

var _ resourceids.ResourceId = &RunId{}

// RunId is a struct representing the Resource ID for a Run
type RunId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	RunId             string
}

// NewRunID returns a new RunId struct
func NewRunID(subscriptionId string, resourceGroupName string, registryName string, runId string) RunId {
	return RunId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		RunId:             runId,
	}
}

// ParseRunID parses 'input' into a RunId
func ParseRunID(input string) (*RunId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RunId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RunId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRunIDInsensitively parses 'input' case-insensitively into a RunId
// note: this method should only be used for API response data and not user input
func ParseRunIDInsensitively(input string) (*RunId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RunId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RunId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RunId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.RunId, ok = input.Parsed["runId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "runId", input)
	}

	return nil
}

// ValidateRunID checks that 'input' can be parsed as a Run ID
func ValidateRunID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRunID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Run ID
func (id RunId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/runs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.RunId)
}

// Segments returns a slice of Resource ID Segments which comprise this Run ID
func (id RunId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryName"),
		resourceids.StaticSegment("staticRuns", "runs", "runs"),
		resourceids.UserSpecifiedSegment("runId", "runId"),
	}
}

// String returns a human-readable description of this Run ID
func (id RunId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Run: %q", id.RunId),
	}
	return fmt.Sprintf("Run (%s)", strings.Join(components, "\n"))
}
