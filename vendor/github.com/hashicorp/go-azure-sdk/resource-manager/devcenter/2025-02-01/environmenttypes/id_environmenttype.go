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
	recaser.RegisterResourceId(&EnvironmentTypeId{})
}

var _ resourceids.ResourceId = &EnvironmentTypeId{}

// EnvironmentTypeId is a struct representing the Resource ID for a Environment Type
type EnvironmentTypeId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ProjectName         string
	EnvironmentTypeName string
}

// NewEnvironmentTypeID returns a new EnvironmentTypeId struct
func NewEnvironmentTypeID(subscriptionId string, resourceGroupName string, projectName string, environmentTypeName string) EnvironmentTypeId {
	return EnvironmentTypeId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ProjectName:         projectName,
		EnvironmentTypeName: environmentTypeName,
	}
}

// ParseEnvironmentTypeID parses 'input' into a EnvironmentTypeId
func ParseEnvironmentTypeID(input string) (*EnvironmentTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EnvironmentTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EnvironmentTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseEnvironmentTypeIDInsensitively parses 'input' case-insensitively into a EnvironmentTypeId
// note: this method should only be used for API response data and not user input
func ParseEnvironmentTypeIDInsensitively(input string) (*EnvironmentTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EnvironmentTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EnvironmentTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *EnvironmentTypeId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.EnvironmentTypeName, ok = input.Parsed["environmentTypeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "environmentTypeName", input)
	}

	return nil
}

// ValidateEnvironmentTypeID checks that 'input' can be parsed as a Environment Type ID
func ValidateEnvironmentTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEnvironmentTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Environment Type ID
func (id EnvironmentTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/projects/%s/environmentTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProjectName, id.EnvironmentTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Environment Type ID
func (id EnvironmentTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticProjects", "projects", "projects"),
		resourceids.UserSpecifiedSegment("projectName", "projectName"),
		resourceids.StaticSegment("staticEnvironmentTypes", "environmentTypes", "environmentTypes"),
		resourceids.UserSpecifiedSegment("environmentTypeName", "environmentTypeName"),
	}
}

// String returns a human-readable description of this Environment Type ID
func (id EnvironmentTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Environment Type Name: %q", id.EnvironmentTypeName),
	}
	return fmt.Sprintf("Environment Type (%s)", strings.Join(components, "\n"))
}
