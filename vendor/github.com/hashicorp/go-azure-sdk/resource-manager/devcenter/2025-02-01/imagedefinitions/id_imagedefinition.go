package imagedefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ImageDefinitionId{})
}

var _ resourceids.ResourceId = &ImageDefinitionId{}

// ImageDefinitionId is a struct representing the Resource ID for a Image Definition
type ImageDefinitionId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ProjectName         string
	CatalogName         string
	ImageDefinitionName string
}

// NewImageDefinitionID returns a new ImageDefinitionId struct
func NewImageDefinitionID(subscriptionId string, resourceGroupName string, projectName string, catalogName string, imageDefinitionName string) ImageDefinitionId {
	return ImageDefinitionId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ProjectName:         projectName,
		CatalogName:         catalogName,
		ImageDefinitionName: imageDefinitionName,
	}
}

// ParseImageDefinitionID parses 'input' into a ImageDefinitionId
func ParseImageDefinitionID(input string) (*ImageDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ImageDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ImageDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseImageDefinitionIDInsensitively parses 'input' case-insensitively into a ImageDefinitionId
// note: this method should only be used for API response data and not user input
func ParseImageDefinitionIDInsensitively(input string) (*ImageDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ImageDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ImageDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ImageDefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ProjectName, ok = input.Parsed["projectName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "projectName", input)
	}

	if id.CatalogName, ok = input.Parsed["catalogName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "catalogName", input)
	}

	if id.ImageDefinitionName, ok = input.Parsed["imageDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "imageDefinitionName", input)
	}

	return nil
}

// ValidateImageDefinitionID checks that 'input' can be parsed as a Image Definition ID
func ValidateImageDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseImageDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Image Definition ID
func (id ImageDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/projects/%s/catalogs/%s/imageDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProjectName, id.CatalogName, id.ImageDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Image Definition ID
func (id ImageDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticProjects", "projects", "projects"),
		resourceids.UserSpecifiedSegment("projectName", "projectName"),
		resourceids.StaticSegment("staticCatalogs", "catalogs", "catalogs"),
		resourceids.UserSpecifiedSegment("catalogName", "catalogName"),
		resourceids.StaticSegment("staticImageDefinitions", "imageDefinitions", "imageDefinitions"),
		resourceids.UserSpecifiedSegment("imageDefinitionName", "imageDefinitionName"),
	}
}

// String returns a human-readable description of this Image Definition ID
func (id ImageDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Catalog Name: %q", id.CatalogName),
		fmt.Sprintf("Image Definition Name: %q", id.ImageDefinitionName),
	}
	return fmt.Sprintf("Image Definition (%s)", strings.Join(components, "\n"))
}
