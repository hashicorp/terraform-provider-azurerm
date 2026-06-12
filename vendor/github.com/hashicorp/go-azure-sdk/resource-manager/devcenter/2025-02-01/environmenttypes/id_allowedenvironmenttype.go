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
	recaser.RegisterResourceId(&AllowedEnvironmentTypeId{})
}

var _ resourceids.ResourceId = &AllowedEnvironmentTypeId{}

// AllowedEnvironmentTypeId is a struct representing the Resource ID for a Allowed Environment Type
type AllowedEnvironmentTypeId struct {
	SubscriptionId             string
	ResourceGroupName          string
	ProjectName                string
	AllowedEnvironmentTypeName string
}

// NewAllowedEnvironmentTypeID returns a new AllowedEnvironmentTypeId struct
func NewAllowedEnvironmentTypeID(subscriptionId string, resourceGroupName string, projectName string, allowedEnvironmentTypeName string) AllowedEnvironmentTypeId {
	return AllowedEnvironmentTypeId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		ProjectName:                projectName,
		AllowedEnvironmentTypeName: allowedEnvironmentTypeName,
	}
}

// ParseAllowedEnvironmentTypeID parses 'input' into a AllowedEnvironmentTypeId
func ParseAllowedEnvironmentTypeID(input string) (*AllowedEnvironmentTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AllowedEnvironmentTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AllowedEnvironmentTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAllowedEnvironmentTypeIDInsensitively parses 'input' case-insensitively into a AllowedEnvironmentTypeId
// note: this method should only be used for API response data and not user input
func ParseAllowedEnvironmentTypeIDInsensitively(input string) (*AllowedEnvironmentTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AllowedEnvironmentTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AllowedEnvironmentTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AllowedEnvironmentTypeId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AllowedEnvironmentTypeName, ok = input.Parsed["allowedEnvironmentTypeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "allowedEnvironmentTypeName", input)
	}

	return nil
}

// ValidateAllowedEnvironmentTypeID checks that 'input' can be parsed as a Allowed Environment Type ID
func ValidateAllowedEnvironmentTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAllowedEnvironmentTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Allowed Environment Type ID
func (id AllowedEnvironmentTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/projects/%s/allowedEnvironmentTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProjectName, id.AllowedEnvironmentTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Allowed Environment Type ID
func (id AllowedEnvironmentTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticProjects", "projects", "projects"),
		resourceids.UserSpecifiedSegment("projectName", "projectName"),
		resourceids.StaticSegment("staticAllowedEnvironmentTypes", "allowedEnvironmentTypes", "allowedEnvironmentTypes"),
		resourceids.UserSpecifiedSegment("allowedEnvironmentTypeName", "allowedEnvironmentTypeName"),
	}
}

// String returns a human-readable description of this Allowed Environment Type ID
func (id AllowedEnvironmentTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Allowed Environment Type Name: %q", id.AllowedEnvironmentTypeName),
	}
	return fmt.Sprintf("Allowed Environment Type (%s)", strings.Join(components, "\n"))
}
