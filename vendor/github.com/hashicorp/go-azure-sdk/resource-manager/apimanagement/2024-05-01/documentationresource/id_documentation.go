package documentationresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DocumentationId{})
}

var _ resourceids.ResourceId = &DocumentationId{}

// DocumentationId is a struct representing the Resource ID for a Documentation
type DocumentationId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	DocumentationId   string
}

// NewDocumentationID returns a new DocumentationId struct
func NewDocumentationID(subscriptionId string, resourceGroupName string, serviceName string, documentationId string) DocumentationId {
	return DocumentationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		DocumentationId:   documentationId,
	}
}

// ParseDocumentationID parses 'input' into a DocumentationId
func ParseDocumentationID(input string) (*DocumentationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DocumentationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DocumentationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDocumentationIDInsensitively parses 'input' case-insensitively into a DocumentationId
// note: this method should only be used for API response data and not user input
func ParseDocumentationIDInsensitively(input string) (*DocumentationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DocumentationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DocumentationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DocumentationId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DocumentationId, ok = input.Parsed["documentationId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "documentationId", input)
	}

	return nil
}

// ValidateDocumentationID checks that 'input' can be parsed as a Documentation ID
func ValidateDocumentationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDocumentationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Documentation ID
func (id DocumentationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/documentations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.DocumentationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Documentation ID
func (id DocumentationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticDocumentations", "documentations", "documentations"),
		resourceids.UserSpecifiedSegment("documentationId", "documentationId"),
	}
}

// String returns a human-readable description of this Documentation ID
func (id DocumentationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Documentation: %q", id.DocumentationId),
	}
	return fmt.Sprintf("Documentation (%s)", strings.Join(components, "\n"))
}
