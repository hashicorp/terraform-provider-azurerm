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
	recaser.RegisterResourceId(&BuildId{})
}

var _ resourceids.ResourceId = &BuildId{}

// BuildId is a struct representing the Resource ID for a Build
type BuildId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ProjectName         string
	CatalogName         string
	ImageDefinitionName string
	BuildName           string
}

// NewBuildID returns a new BuildId struct
func NewBuildID(subscriptionId string, resourceGroupName string, projectName string, catalogName string, imageDefinitionName string, buildName string) BuildId {
	return BuildId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ProjectName:         projectName,
		CatalogName:         catalogName,
		ImageDefinitionName: imageDefinitionName,
		BuildName:           buildName,
	}
}

// ParseBuildID parses 'input' into a BuildId
func ParseBuildID(input string) (*BuildId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuildId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuildId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBuildIDInsensitively parses 'input' case-insensitively into a BuildId
// note: this method should only be used for API response data and not user input
func ParseBuildIDInsensitively(input string) (*BuildId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuildId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuildId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BuildId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.BuildName, ok = input.Parsed["buildName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "buildName", input)
	}

	return nil
}

// ValidateBuildID checks that 'input' can be parsed as a Build ID
func ValidateBuildID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBuildID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Build ID
func (id BuildId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/projects/%s/catalogs/%s/imageDefinitions/%s/builds/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProjectName, id.CatalogName, id.ImageDefinitionName, id.BuildName)
}

// Segments returns a slice of Resource ID Segments which comprise this Build ID
func (id BuildId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticBuilds", "builds", "builds"),
		resourceids.UserSpecifiedSegment("buildName", "buildName"),
	}
}

// String returns a human-readable description of this Build ID
func (id BuildId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Catalog Name: %q", id.CatalogName),
		fmt.Sprintf("Image Definition Name: %q", id.ImageDefinitionName),
		fmt.Sprintf("Build Name: %q", id.BuildName),
	}
	return fmt.Sprintf("Build (%s)", strings.Join(components, "\n"))
}
