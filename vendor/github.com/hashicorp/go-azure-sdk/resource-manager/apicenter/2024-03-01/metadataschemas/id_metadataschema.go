package metadataschemas

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MetadataSchemaId{})
}

var _ resourceids.ResourceId = &MetadataSchemaId{}

// MetadataSchemaId is a struct representing the Resource ID for a Metadata Schema
type MetadataSchemaId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ServiceName        string
	MetadataSchemaName string
}

// NewMetadataSchemaID returns a new MetadataSchemaId struct
func NewMetadataSchemaID(subscriptionId string, resourceGroupName string, serviceName string, metadataSchemaName string) MetadataSchemaId {
	return MetadataSchemaId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ServiceName:        serviceName,
		MetadataSchemaName: metadataSchemaName,
	}
}

// ParseMetadataSchemaID parses 'input' into a MetadataSchemaId
func ParseMetadataSchemaID(input string) (*MetadataSchemaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MetadataSchemaId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MetadataSchemaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMetadataSchemaIDInsensitively parses 'input' case-insensitively into a MetadataSchemaId
// note: this method should only be used for API response data and not user input
func ParseMetadataSchemaIDInsensitively(input string) (*MetadataSchemaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MetadataSchemaId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MetadataSchemaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MetadataSchemaId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.MetadataSchemaName, ok = input.Parsed["metadataSchemaName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "metadataSchemaName", input)
	}

	return nil
}

// ValidateMetadataSchemaID checks that 'input' can be parsed as a Metadata Schema ID
func ValidateMetadataSchemaID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMetadataSchemaID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Metadata Schema ID
func (id MetadataSchemaId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiCenter/services/%s/metadataSchemas/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.MetadataSchemaName)
}

// Segments returns a slice of Resource ID Segments which comprise this Metadata Schema ID
func (id MetadataSchemaId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiCenter", "Microsoft.ApiCenter", "Microsoft.ApiCenter"),
		resourceids.StaticSegment("staticServices", "services", "services"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticMetadataSchemas", "metadataSchemas", "metadataSchemas"),
		resourceids.UserSpecifiedSegment("metadataSchemaName", "metadataSchemaName"),
	}
}

// String returns a human-readable description of this Metadata Schema ID
func (id MetadataSchemaId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Metadata Schema Name: %q", id.MetadataSchemaName),
	}
	return fmt.Sprintf("Metadata Schema (%s)", strings.Join(components, "\n"))
}
