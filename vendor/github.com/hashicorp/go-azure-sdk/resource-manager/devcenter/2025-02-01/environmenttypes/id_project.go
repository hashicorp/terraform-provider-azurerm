package environmenttypes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProjectId{})
}

var _ resourceids.ResourceId = &ProjectId{}

// ProjectId is a struct representing the Resource ID for a Project
type ProjectId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProjectName       string
}

// NewProjectID returns a new ProjectId struct
func NewProjectID(subscriptionId string, resourceGroupName string, projectName string) ProjectId {
	return ProjectId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProjectName:       projectName,
	}
}

// ParseProjectID parses 'input' into a ProjectId
func ParseProjectID(input string) (*ProjectId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProjectId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProjectId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProjectIDInsensitively parses 'input' case-insensitively into a ProjectId
// note: this method should only be used for API response data and not user input
func ParseProjectIDInsensitively(input string) (*ProjectId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProjectId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProjectId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProjectId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateProjectID checks that 'input' can be parsed as a Project ID
func ValidateProjectID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProjectID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Project ID
func (id ProjectId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/projects/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProjectName)
}

// Segments returns a slice of Resource ID Segments which comprise this Project ID
func (id ProjectId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticProjects", "projects", "projects"),
		resourceids.UserSpecifiedSegment("projectName", "projectName"),
	}
}

// String returns a human-readable description of this Project ID
func (id ProjectId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
	}
	return fmt.Sprintf("Project (%s)", strings.Join(components, "\n"))
}
