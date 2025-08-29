package apitagdescription

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TagDescriptionId{})
}

var _ resourceids.ResourceId = &TagDescriptionId{}

// TagDescriptionId is a struct representing the Resource ID for a Tag Description
type TagDescriptionId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ApiId             string
	TagDescriptionId  string
}

// NewTagDescriptionID returns a new TagDescriptionId struct
func NewTagDescriptionID(subscriptionId string, resourceGroupName string, serviceName string, apiId string, tagDescriptionId string) TagDescriptionId {
	return TagDescriptionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ApiId:             apiId,
		TagDescriptionId:  tagDescriptionId,
	}
}

// ParseTagDescriptionID parses 'input' into a TagDescriptionId
func ParseTagDescriptionID(input string) (*TagDescriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TagDescriptionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TagDescriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTagDescriptionIDInsensitively parses 'input' case-insensitively into a TagDescriptionId
// note: this method should only be used for API response data and not user input
func ParseTagDescriptionIDInsensitively(input string) (*TagDescriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TagDescriptionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TagDescriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TagDescriptionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.TagDescriptionId, ok = input.Parsed["tagDescriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tagDescriptionId", input)
	}

	return nil
}

// ValidateTagDescriptionID checks that 'input' can be parsed as a Tag Description ID
func ValidateTagDescriptionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTagDescriptionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Tag Description ID
func (id TagDescriptionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/tagDescriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.TagDescriptionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Tag Description ID
func (id TagDescriptionId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticTagDescriptions", "tagDescriptions", "tagDescriptions"),
		resourceids.UserSpecifiedSegment("tagDescriptionId", "tagDescriptionId"),
	}
}

// String returns a human-readable description of this Tag Description ID
func (id TagDescriptionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Api: %q", id.ApiId),
		fmt.Sprintf("Tag Description: %q", id.TagDescriptionId),
	}
	return fmt.Sprintf("Tag Description (%s)", strings.Join(components, "\n"))
}
