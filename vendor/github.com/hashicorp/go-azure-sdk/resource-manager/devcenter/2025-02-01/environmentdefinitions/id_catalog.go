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
	recaser.RegisterResourceId(&CatalogId{})
}

var _ resourceids.ResourceId = &CatalogId{}

// CatalogId is a struct representing the Resource ID for a Catalog
type CatalogId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProjectName       string
	CatalogName       string
}

// NewCatalogID returns a new CatalogId struct
func NewCatalogID(subscriptionId string, resourceGroupName string, projectName string, catalogName string) CatalogId {
	return CatalogId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProjectName:       projectName,
		CatalogName:       catalogName,
	}
}

// ParseCatalogID parses 'input' into a CatalogId
func ParseCatalogID(input string) (*CatalogId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CatalogId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CatalogId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCatalogIDInsensitively parses 'input' case-insensitively into a CatalogId
// note: this method should only be used for API response data and not user input
func ParseCatalogIDInsensitively(input string) (*CatalogId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CatalogId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CatalogId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CatalogId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateCatalogID checks that 'input' can be parsed as a Catalog ID
func ValidateCatalogID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCatalogID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Catalog ID
func (id CatalogId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/projects/%s/catalogs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProjectName, id.CatalogName)
}

// Segments returns a slice of Resource ID Segments which comprise this Catalog ID
func (id CatalogId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Catalog ID
func (id CatalogId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Catalog Name: %q", id.CatalogName),
	}
	return fmt.Sprintf("Catalog (%s)", strings.Join(components, "\n"))
}
