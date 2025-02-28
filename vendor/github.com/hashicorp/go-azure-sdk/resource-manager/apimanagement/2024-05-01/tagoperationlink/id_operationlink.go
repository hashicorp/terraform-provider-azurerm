package tagoperationlink

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OperationLinkId{})
}

var _ resourceids.ResourceId = &OperationLinkId{}

// OperationLinkId is a struct representing the Resource ID for a Operation Link
type OperationLinkId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	TagId             string
	OperationLinkId   string
}

// NewOperationLinkID returns a new OperationLinkId struct
func NewOperationLinkID(subscriptionId string, resourceGroupName string, serviceName string, tagId string, operationLinkId string) OperationLinkId {
	return OperationLinkId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		TagId:             tagId,
		OperationLinkId:   operationLinkId,
	}
}

// ParseOperationLinkID parses 'input' into a OperationLinkId
func ParseOperationLinkID(input string) (*OperationLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOperationLinkIDInsensitively parses 'input' case-insensitively into a OperationLinkId
// note: this method should only be used for API response data and not user input
func ParseOperationLinkIDInsensitively(input string) (*OperationLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OperationLinkId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.TagId, ok = input.Parsed["tagId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tagId", input)
	}

	if id.OperationLinkId, ok = input.Parsed["operationLinkId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "operationLinkId", input)
	}

	return nil
}

// ValidateOperationLinkID checks that 'input' can be parsed as a Operation Link ID
func ValidateOperationLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOperationLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Operation Link ID
func (id OperationLinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/tags/%s/operationLinks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.TagId, id.OperationLinkId)
}

// Segments returns a slice of Resource ID Segments which comprise this Operation Link ID
func (id OperationLinkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticTags", "tags", "tags"),
		resourceids.UserSpecifiedSegment("tagId", "tagId"),
		resourceids.StaticSegment("staticOperationLinks", "operationLinks", "operationLinks"),
		resourceids.UserSpecifiedSegment("operationLinkId", "operationLinkId"),
	}
}

// String returns a human-readable description of this Operation Link ID
func (id OperationLinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Tag: %q", id.TagId),
		fmt.Sprintf("Operation Link: %q", id.OperationLinkId),
	}
	return fmt.Sprintf("Operation Link (%s)", strings.Join(components, "\n"))
}
