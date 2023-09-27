package importpipelines

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ImportPipelineId{}

// ImportPipelineId is a struct representing the Resource ID for a Import Pipeline
type ImportPipelineId struct {
	SubscriptionId     string
	ResourceGroupName  string
	RegistryName       string
	ImportPipelineName string
}

// NewImportPipelineID returns a new ImportPipelineId struct
func NewImportPipelineID(subscriptionId string, resourceGroupName string, registryName string, importPipelineName string) ImportPipelineId {
	return ImportPipelineId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		RegistryName:       registryName,
		ImportPipelineName: importPipelineName,
	}
}

// ParseImportPipelineID parses 'input' into a ImportPipelineId
func ParseImportPipelineID(input string) (*ImportPipelineId, error) {
	parser := resourceids.NewParserFromResourceIdType(ImportPipelineId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ImportPipelineId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.ImportPipelineName, ok = parsed.Parsed["importPipelineName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "importPipelineName", *parsed)
	}

	return &id, nil
}

// ParseImportPipelineIDInsensitively parses 'input' case-insensitively into a ImportPipelineId
// note: this method should only be used for API response data and not user input
func ParseImportPipelineIDInsensitively(input string) (*ImportPipelineId, error) {
	parser := resourceids.NewParserFromResourceIdType(ImportPipelineId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ImportPipelineId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.ImportPipelineName, ok = parsed.Parsed["importPipelineName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "importPipelineName", *parsed)
	}

	return &id, nil
}

// ValidateImportPipelineID checks that 'input' can be parsed as a Import Pipeline ID
func ValidateImportPipelineID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseImportPipelineID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Import Pipeline ID
func (id ImportPipelineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/importPipelines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.ImportPipelineName)
}

// Segments returns a slice of Resource ID Segments which comprise this Import Pipeline ID
func (id ImportPipelineId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryValue"),
		resourceids.StaticSegment("staticImportPipelines", "importPipelines", "importPipelines"),
		resourceids.UserSpecifiedSegment("importPipelineName", "importPipelineValue"),
	}
}

// String returns a human-readable description of this Import Pipeline ID
func (id ImportPipelineId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Import Pipeline Name: %q", id.ImportPipelineName),
	}
	return fmt.Sprintf("Import Pipeline (%s)", strings.Join(components, "\n"))
}
