package tagapilink

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TagApiLinkId{})
}

var _ resourceids.ResourceId = &TagApiLinkId{}

// TagApiLinkId is a struct representing the Resource ID for a Tag Api Link
type TagApiLinkId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	TagId             string
	ApiLinkId         string
}

// NewTagApiLinkID returns a new TagApiLinkId struct
func NewTagApiLinkID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, tagId string, apiLinkId string) TagApiLinkId {
	return TagApiLinkId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		TagId:             tagId,
		ApiLinkId:         apiLinkId,
	}
}

// ParseTagApiLinkID parses 'input' into a TagApiLinkId
func ParseTagApiLinkID(input string) (*TagApiLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TagApiLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TagApiLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTagApiLinkIDInsensitively parses 'input' case-insensitively into a TagApiLinkId
// note: this method should only be used for API response data and not user input
func ParseTagApiLinkIDInsensitively(input string) (*TagApiLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TagApiLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TagApiLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TagApiLinkId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.WorkspaceId, ok = input.Parsed["workspaceId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceId", input)
	}

	if id.TagId, ok = input.Parsed["tagId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tagId", input)
	}

	if id.ApiLinkId, ok = input.Parsed["apiLinkId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiLinkId", input)
	}

	return nil
}

// ValidateTagApiLinkID checks that 'input' can be parsed as a Tag Api Link ID
func ValidateTagApiLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTagApiLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Tag Api Link ID
func (id TagApiLinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/tags/%s/apiLinks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.TagId, id.ApiLinkId)
}

// Segments returns a slice of Resource ID Segments which comprise this Tag Api Link ID
func (id TagApiLinkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceId", "workspaceId"),
		resourceids.StaticSegment("staticTags", "tags", "tags"),
		resourceids.UserSpecifiedSegment("tagId", "tagId"),
		resourceids.StaticSegment("staticApiLinks", "apiLinks", "apiLinks"),
		resourceids.UserSpecifiedSegment("apiLinkId", "apiLinkId"),
	}
}

// String returns a human-readable description of this Tag Api Link ID
func (id TagApiLinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Tag: %q", id.TagId),
		fmt.Sprintf("Api Link: %q", id.ApiLinkId),
	}
	return fmt.Sprintf("Tag Api Link (%s)", strings.Join(components, "\n"))
}
