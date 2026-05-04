package apioperationtag

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OperationTagId{})
}

var _ resourceids.ResourceId = &OperationTagId{}

// OperationTagId is a struct representing the Resource ID for a Operation Tag
type OperationTagId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ApiId             string
	OperationId       string
	TagId             string
}

// NewOperationTagID returns a new OperationTagId struct
func NewOperationTagID(subscriptionId string, resourceGroupName string, serviceName string, apiId string, operationId string, tagId string) OperationTagId {
	return OperationTagId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ApiId:             apiId,
		OperationId:       operationId,
		TagId:             tagId,
	}
}

// ParseOperationTagID parses 'input' into a OperationTagId
func ParseOperationTagID(input string) (*OperationTagId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationTagId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationTagId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOperationTagIDInsensitively parses 'input' case-insensitively into a OperationTagId
// note: this method should only be used for API response data and not user input
func ParseOperationTagIDInsensitively(input string) (*OperationTagId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationTagId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationTagId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OperationTagId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.OperationId, ok = input.Parsed["operationId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "operationId", input)
	}

	if id.TagId, ok = input.Parsed["tagId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tagId", input)
	}

	return nil
}

// ValidateOperationTagID checks that 'input' can be parsed as a Operation Tag ID
func ValidateOperationTagID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOperationTagID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Operation Tag ID
func (id OperationTagId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/operations/%s/tags/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.OperationId, id.TagId)
}

// Segments returns a slice of Resource ID Segments which comprise this Operation Tag ID
func (id OperationTagId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticOperations", "operations", "operations"),
		resourceids.UserSpecifiedSegment("operationId", "operationId"),
		resourceids.StaticSegment("staticTags", "tags", "tags"),
		resourceids.UserSpecifiedSegment("tagId", "tagId"),
	}
}

// String returns a human-readable description of this Operation Tag ID
func (id OperationTagId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Api: %q", id.ApiId),
		fmt.Sprintf("Operation: %q", id.OperationId),
		fmt.Sprintf("Tag: %q", id.TagId),
	}
	return fmt.Sprintf("Operation Tag (%s)", strings.Join(components, "\n"))
}
