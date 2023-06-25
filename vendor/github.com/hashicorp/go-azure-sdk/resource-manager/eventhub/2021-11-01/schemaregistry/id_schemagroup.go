package schemaregistry

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SchemaGroupId{}

// SchemaGroupId is a struct representing the Resource ID for a Schema Group
type SchemaGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	SchemaGroupName   string
}

// NewSchemaGroupID returns a new SchemaGroupId struct
func NewSchemaGroupID(subscriptionId string, resourceGroupName string, namespaceName string, schemaGroupName string) SchemaGroupId {
	return SchemaGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		SchemaGroupName:   schemaGroupName,
	}
}

// ParseSchemaGroupID parses 'input' into a SchemaGroupId
func ParseSchemaGroupID(input string) (*SchemaGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(SchemaGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SchemaGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.SchemaGroupName, ok = parsed.Parsed["schemaGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "schemaGroupName", *parsed)
	}

	return &id, nil
}

// ParseSchemaGroupIDInsensitively parses 'input' case-insensitively into a SchemaGroupId
// note: this method should only be used for API response data and not user input
func ParseSchemaGroupIDInsensitively(input string) (*SchemaGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(SchemaGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SchemaGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", *parsed)
	}

	if id.SchemaGroupName, ok = parsed.Parsed["schemaGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "schemaGroupName", *parsed)
	}

	return &id, nil
}

// ValidateSchemaGroupID checks that 'input' can be parsed as a Schema Group ID
func ValidateSchemaGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSchemaGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Schema Group ID
func (id SchemaGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/schemaGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.SchemaGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Schema Group ID
func (id SchemaGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventHub", "Microsoft.EventHub", "Microsoft.EventHub"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticSchemaGroups", "schemaGroups", "schemaGroups"),
		resourceids.UserSpecifiedSegment("schemaGroupName", "schemaGroupValue"),
	}
}

// String returns a human-readable description of this Schema Group ID
func (id SchemaGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Schema Group Name: %q", id.SchemaGroupName),
	}
	return fmt.Sprintf("Schema Group (%s)", strings.Join(components, "\n"))
}
