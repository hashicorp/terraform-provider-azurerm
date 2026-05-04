package apischema

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApiSchemaId{})
}

var _ resourceids.ResourceId = &ApiSchemaId{}

// ApiSchemaId is a struct representing the Resource ID for a Api Schema
type ApiSchemaId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ApiId             string
	SchemaId          string
}

// NewApiSchemaID returns a new ApiSchemaId struct
func NewApiSchemaID(subscriptionId string, resourceGroupName string, serviceName string, apiId string, schemaId string) ApiSchemaId {
	return ApiSchemaId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ApiId:             apiId,
		SchemaId:          schemaId,
	}
}

// ParseApiSchemaID parses 'input' into a ApiSchemaId
func ParseApiSchemaID(input string) (*ApiSchemaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiSchemaId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiSchemaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApiSchemaIDInsensitively parses 'input' case-insensitively into a ApiSchemaId
// note: this method should only be used for API response data and not user input
func ParseApiSchemaIDInsensitively(input string) (*ApiSchemaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiSchemaId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiSchemaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApiSchemaId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.ApiId, ok = input.Parsed["apiId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiId", input)
	}

	if id.SchemaId, ok = input.Parsed["schemaId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "schemaId", input)
	}

	return nil
}

// ValidateApiSchemaID checks that 'input' can be parsed as a Api Schema ID
func ValidateApiSchemaID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApiSchemaID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Api Schema ID
func (id ApiSchemaId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/schemas/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.SchemaId)
}

// Segments returns a slice of Resource ID Segments which comprise this Api Schema ID
func (id ApiSchemaId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticApis", "apis", "apis"),
		resourceids.UserSpecifiedSegment("apiId", "apiId"),
		resourceids.StaticSegment("staticSchemas", "schemas", "schemas"),
		resourceids.UserSpecifiedSegment("schemaId", "schemaId"),
	}
}

// String returns a human-readable description of this Api Schema ID
func (id ApiSchemaId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Api: %q", id.ApiId),
		fmt.Sprintf("Schema: %q", id.SchemaId),
	}
	return fmt.Sprintf("Api Schema (%s)", strings.Join(components, "\n"))
}
