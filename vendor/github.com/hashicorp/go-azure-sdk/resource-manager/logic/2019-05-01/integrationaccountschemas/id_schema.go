package integrationaccountschemas

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SchemaId{})
}

var _ resourceids.ResourceId = &SchemaId{}

// SchemaId is a struct representing the Resource ID for a Schema
type SchemaId struct {
	SubscriptionId         string
	ResourceGroupName      string
	IntegrationAccountName string
	SchemaName             string
}

// NewSchemaID returns a new SchemaId struct
func NewSchemaID(subscriptionId string, resourceGroupName string, integrationAccountName string, schemaName string) SchemaId {
	return SchemaId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		IntegrationAccountName: integrationAccountName,
		SchemaName:             schemaName,
	}
}

// ParseSchemaID parses 'input' into a SchemaId
func ParseSchemaID(input string) (*SchemaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SchemaId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SchemaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSchemaIDInsensitively parses 'input' case-insensitively into a SchemaId
// note: this method should only be used for API response data and not user input
func ParseSchemaIDInsensitively(input string) (*SchemaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SchemaId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SchemaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SchemaId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.IntegrationAccountName, ok = input.Parsed["integrationAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "integrationAccountName", input)
	}

	if id.SchemaName, ok = input.Parsed["schemaName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "schemaName", input)
	}

	return nil
}

// ValidateSchemaID checks that 'input' can be parsed as a Schema ID
func ValidateSchemaID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSchemaID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Schema ID
func (id SchemaId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/schemas/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IntegrationAccountName, id.SchemaName)
}

// Segments returns a slice of Resource ID Segments which comprise this Schema ID
func (id SchemaId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticIntegrationAccounts", "integrationAccounts", "integrationAccounts"),
		resourceids.UserSpecifiedSegment("integrationAccountName", "integrationAccountName"),
		resourceids.StaticSegment("staticSchemas", "schemas", "schemas"),
		resourceids.UserSpecifiedSegment("schemaName", "schemaName"),
	}
}

// String returns a human-readable description of this Schema ID
func (id SchemaId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Integration Account Name: %q", id.IntegrationAccountName),
		fmt.Sprintf("Schema Name: %q", id.SchemaName),
	}
	return fmt.Sprintf("Schema (%s)", strings.Join(components, "\n"))
}
