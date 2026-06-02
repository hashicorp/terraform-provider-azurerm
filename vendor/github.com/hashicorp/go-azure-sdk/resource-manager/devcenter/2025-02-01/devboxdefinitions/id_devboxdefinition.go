package devboxdefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DevBoxDefinitionId{})
}

var _ resourceids.ResourceId = &DevBoxDefinitionId{}

// DevBoxDefinitionId is a struct representing the Resource ID for a Dev Box Definition
type DevBoxDefinitionId struct {
	SubscriptionId       string
	ResourceGroupName    string
	ProjectName          string
	DevBoxDefinitionName string
}

// NewDevBoxDefinitionID returns a new DevBoxDefinitionId struct
func NewDevBoxDefinitionID(subscriptionId string, resourceGroupName string, projectName string, devBoxDefinitionName string) DevBoxDefinitionId {
	return DevBoxDefinitionId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		ProjectName:          projectName,
		DevBoxDefinitionName: devBoxDefinitionName,
	}
}

// ParseDevBoxDefinitionID parses 'input' into a DevBoxDefinitionId
func ParseDevBoxDefinitionID(input string) (*DevBoxDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DevBoxDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DevBoxDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDevBoxDefinitionIDInsensitively parses 'input' case-insensitively into a DevBoxDefinitionId
// note: this method should only be used for API response data and not user input
func ParseDevBoxDefinitionIDInsensitively(input string) (*DevBoxDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DevBoxDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DevBoxDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DevBoxDefinitionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DevBoxDefinitionName, ok = input.Parsed["devBoxDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "devBoxDefinitionName", input)
	}

	return nil
}

// ValidateDevBoxDefinitionID checks that 'input' can be parsed as a Dev Box Definition ID
func ValidateDevBoxDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDevBoxDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dev Box Definition ID
func (id DevBoxDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/projects/%s/devBoxDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProjectName, id.DevBoxDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dev Box Definition ID
func (id DevBoxDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticProjects", "projects", "projects"),
		resourceids.UserSpecifiedSegment("projectName", "projectName"),
		resourceids.StaticSegment("staticDevBoxDefinitions", "devBoxDefinitions", "devBoxDefinitions"),
		resourceids.UserSpecifiedSegment("devBoxDefinitionName", "devBoxDefinitionName"),
	}
}

// String returns a human-readable description of this Dev Box Definition ID
func (id DevBoxDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Dev Box Definition Name: %q", id.DevBoxDefinitionName),
	}
	return fmt.Sprintf("Dev Box Definition (%s)", strings.Join(components, "\n"))
}
