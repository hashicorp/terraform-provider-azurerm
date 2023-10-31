package apitag

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ApiTagId{}

// ApiTagId is a struct representing the Resource ID for a Api Tag
type ApiTagId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ApiId             string
	TagId             string
}

// NewApiTagID returns a new ApiTagId struct
func NewApiTagID(subscriptionId string, resourceGroupName string, serviceName string, apiId string, tagId string) ApiTagId {
	return ApiTagId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ApiId:             apiId,
		TagId:             tagId,
	}
}

// ParseApiTagID parses 'input' into a ApiTagId
func ParseApiTagID(input string) (*ApiTagId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApiTagId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApiTagId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.ApiId, ok = parsed.Parsed["apiId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "apiId", *parsed)
	}

	if id.TagId, ok = parsed.Parsed["tagId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "tagId", *parsed)
	}

	return &id, nil
}

// ParseApiTagIDInsensitively parses 'input' case-insensitively into a ApiTagId
// note: this method should only be used for API response data and not user input
func ParseApiTagIDInsensitively(input string) (*ApiTagId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApiTagId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApiTagId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.ApiId, ok = parsed.Parsed["apiId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "apiId", *parsed)
	}

	if id.TagId, ok = parsed.Parsed["tagId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "tagId", *parsed)
	}

	return &id, nil
}

// ValidateApiTagID checks that 'input' can be parsed as a Api Tag ID
func ValidateApiTagID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApiTagID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Api Tag ID
func (id ApiTagId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/tags/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.TagId)
}

// Segments returns a slice of Resource ID Segments which comprise this Api Tag ID
func (id ApiTagId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticApis", "apis", "apis"),
		resourceids.UserSpecifiedSegment("apiId", "apiIdValue"),
		resourceids.StaticSegment("staticTags", "tags", "tags"),
		resourceids.UserSpecifiedSegment("tagId", "tagIdValue"),
	}
}

// String returns a human-readable description of this Api Tag ID
func (id ApiTagId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Api: %q", id.ApiId),
		fmt.Sprintf("Tag: %q", id.TagId),
	}
	return fmt.Sprintf("Api Tag (%s)", strings.Join(components, "\n"))
}
