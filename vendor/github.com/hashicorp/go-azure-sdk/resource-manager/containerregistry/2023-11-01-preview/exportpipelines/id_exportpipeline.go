package exportpipelines

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExportPipelineId{})
}

var _ resourceids.ResourceId = &ExportPipelineId{}

// ExportPipelineId is a struct representing the Resource ID for a Export Pipeline
type ExportPipelineId struct {
	SubscriptionId     string
	ResourceGroupName  string
	RegistryName       string
	ExportPipelineName string
}

// NewExportPipelineID returns a new ExportPipelineId struct
func NewExportPipelineID(subscriptionId string, resourceGroupName string, registryName string, exportPipelineName string) ExportPipelineId {
	return ExportPipelineId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		RegistryName:       registryName,
		ExportPipelineName: exportPipelineName,
	}
}

// ParseExportPipelineID parses 'input' into a ExportPipelineId
func ParseExportPipelineID(input string) (*ExportPipelineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExportPipelineId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExportPipelineId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExportPipelineIDInsensitively parses 'input' case-insensitively into a ExportPipelineId
// note: this method should only be used for API response data and not user input
func ParseExportPipelineIDInsensitively(input string) (*ExportPipelineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExportPipelineId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExportPipelineId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExportPipelineId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ExportPipelineName, ok = input.Parsed["exportPipelineName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "exportPipelineName", input)
	}

	return nil
}

// ValidateExportPipelineID checks that 'input' can be parsed as a Export Pipeline ID
func ValidateExportPipelineID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExportPipelineID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Export Pipeline ID
func (id ExportPipelineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/exportPipelines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.ExportPipelineName)
}

// Segments returns a slice of Resource ID Segments which comprise this Export Pipeline ID
func (id ExportPipelineId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryName"),
		resourceids.StaticSegment("staticExportPipelines", "exportPipelines", "exportPipelines"),
		resourceids.UserSpecifiedSegment("exportPipelineName", "exportPipelineName"),
	}
}

// String returns a human-readable description of this Export Pipeline ID
func (id ExportPipelineId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Export Pipeline Name: %q", id.ExportPipelineName),
	}
	return fmt.Sprintf("Export Pipeline (%s)", strings.Join(components, "\n"))
}
