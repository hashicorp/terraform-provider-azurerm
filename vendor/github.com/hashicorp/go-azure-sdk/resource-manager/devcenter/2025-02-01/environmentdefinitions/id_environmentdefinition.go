package environmentdefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&EnvironmentDefinitionId{})
}

var _ resourceids.ResourceId = &EnvironmentDefinitionId{}

// EnvironmentDefinitionId is a struct representing the Resource ID for a Environment Definition
type EnvironmentDefinitionId struct {
	SubscriptionId            string
	ResourceGroupName         string
	ProjectName               string
	CatalogName               string
	EnvironmentDefinitionName string
}

// NewEnvironmentDefinitionID returns a new EnvironmentDefinitionId struct
func NewEnvironmentDefinitionID(subscriptionId string, resourceGroupName string, projectName string, catalogName string, environmentDefinitionName string) EnvironmentDefinitionId {
	return EnvironmentDefinitionId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		ProjectName:               projectName,
		CatalogName:               catalogName,
		EnvironmentDefinitionName: environmentDefinitionName,
	}
}

// ParseEnvironmentDefinitionID parses 'input' into a EnvironmentDefinitionId
func ParseEnvironmentDefinitionID(input string) (*EnvironmentDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EnvironmentDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EnvironmentDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseEnvironmentDefinitionIDInsensitively parses 'input' case-insensitively into a EnvironmentDefinitionId
// note: this method should only be used for API response data and not user input
func ParseEnvironmentDefinitionIDInsensitively(input string) (*EnvironmentDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EnvironmentDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EnvironmentDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *EnvironmentDefinitionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.EnvironmentDefinitionName, ok = input.Parsed["environmentDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "environmentDefinitionName", input)
	}

	return nil
}

// ValidateEnvironmentDefinitionID checks that 'input' can be parsed as a Environment Definition ID
func ValidateEnvironmentDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEnvironmentDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Environment Definition ID
func (id EnvironmentDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/projects/%s/catalogs/%s/environmentDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProjectName, id.CatalogName, id.EnvironmentDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Environment Definition ID
func (id EnvironmentDefinitionId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticEnvironmentDefinitions", "environmentDefinitions", "environmentDefinitions"),
		resourceids.UserSpecifiedSegment("environmentDefinitionName", "environmentDefinitionName"),
	}
}

// String returns a human-readable description of this Environment Definition ID
func (id EnvironmentDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Catalog Name: %q", id.CatalogName),
		fmt.Sprintf("Environment Definition Name: %q", id.EnvironmentDefinitionName),
	}
	return fmt.Sprintf("Environment Definition (%s)", strings.Join(components, "\n"))
}
