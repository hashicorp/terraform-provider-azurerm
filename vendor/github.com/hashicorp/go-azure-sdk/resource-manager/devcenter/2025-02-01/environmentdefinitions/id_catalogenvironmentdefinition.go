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
	recaser.RegisterResourceId(&CatalogEnvironmentDefinitionId{})
}

var _ resourceids.ResourceId = &CatalogEnvironmentDefinitionId{}

// CatalogEnvironmentDefinitionId is a struct representing the Resource ID for a Catalog Environment Definition
type CatalogEnvironmentDefinitionId struct {
	SubscriptionId            string
	ResourceGroupName         string
	DevCenterName             string
	CatalogName               string
	EnvironmentDefinitionName string
}

// NewCatalogEnvironmentDefinitionID returns a new CatalogEnvironmentDefinitionId struct
func NewCatalogEnvironmentDefinitionID(subscriptionId string, resourceGroupName string, devCenterName string, catalogName string, environmentDefinitionName string) CatalogEnvironmentDefinitionId {
	return CatalogEnvironmentDefinitionId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		DevCenterName:             devCenterName,
		CatalogName:               catalogName,
		EnvironmentDefinitionName: environmentDefinitionName,
	}
}

// ParseCatalogEnvironmentDefinitionID parses 'input' into a CatalogEnvironmentDefinitionId
func ParseCatalogEnvironmentDefinitionID(input string) (*CatalogEnvironmentDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CatalogEnvironmentDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CatalogEnvironmentDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCatalogEnvironmentDefinitionIDInsensitively parses 'input' case-insensitively into a CatalogEnvironmentDefinitionId
// note: this method should only be used for API response data and not user input
func ParseCatalogEnvironmentDefinitionIDInsensitively(input string) (*CatalogEnvironmentDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CatalogEnvironmentDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CatalogEnvironmentDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CatalogEnvironmentDefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DevCenterName, ok = input.Parsed["devCenterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "devCenterName", input)
	}

	if id.CatalogName, ok = input.Parsed["catalogName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "catalogName", input)
	}

	if id.EnvironmentDefinitionName, ok = input.Parsed["environmentDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "environmentDefinitionName", input)
	}

	return nil
}

// ValidateCatalogEnvironmentDefinitionID checks that 'input' can be parsed as a Catalog Environment Definition ID
func ValidateCatalogEnvironmentDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCatalogEnvironmentDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Catalog Environment Definition ID
func (id CatalogEnvironmentDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/devCenters/%s/catalogs/%s/environmentDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DevCenterName, id.CatalogName, id.EnvironmentDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Catalog Environment Definition ID
func (id CatalogEnvironmentDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticDevCenters", "devCenters", "devCenters"),
		resourceids.UserSpecifiedSegment("devCenterName", "devCenterName"),
		resourceids.StaticSegment("staticCatalogs", "catalogs", "catalogs"),
		resourceids.UserSpecifiedSegment("catalogName", "catalogName"),
		resourceids.StaticSegment("staticEnvironmentDefinitions", "environmentDefinitions", "environmentDefinitions"),
		resourceids.UserSpecifiedSegment("environmentDefinitionName", "environmentDefinitionName"),
	}
}

// String returns a human-readable description of this Catalog Environment Definition ID
func (id CatalogEnvironmentDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dev Center Name: %q", id.DevCenterName),
		fmt.Sprintf("Catalog Name: %q", id.CatalogName),
		fmt.Sprintf("Environment Definition Name: %q", id.EnvironmentDefinitionName),
	}
	return fmt.Sprintf("Catalog Environment Definition (%s)", strings.Join(components, "\n"))
}
