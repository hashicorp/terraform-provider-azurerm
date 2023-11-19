package apiversionset

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ApiVersionSetId{}

// ApiVersionSetId is a struct representing the Resource ID for a Api Version Set
type ApiVersionSetId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	VersionSetId      string
}

// NewApiVersionSetID returns a new ApiVersionSetId struct
func NewApiVersionSetID(subscriptionId string, resourceGroupName string, serviceName string, versionSetId string) ApiVersionSetId {
	return ApiVersionSetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		VersionSetId:      versionSetId,
	}
}

// ParseApiVersionSetID parses 'input' into a ApiVersionSetId
func ParseApiVersionSetID(input string) (*ApiVersionSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApiVersionSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApiVersionSetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.VersionSetId, ok = parsed.Parsed["versionSetId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "versionSetId", *parsed)
	}

	return &id, nil
}

// ParseApiVersionSetIDInsensitively parses 'input' case-insensitively into a ApiVersionSetId
// note: this method should only be used for API response data and not user input
func ParseApiVersionSetIDInsensitively(input string) (*ApiVersionSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApiVersionSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApiVersionSetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.VersionSetId, ok = parsed.Parsed["versionSetId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "versionSetId", *parsed)
	}

	return &id, nil
}

// ValidateApiVersionSetID checks that 'input' can be parsed as a Api Version Set ID
func ValidateApiVersionSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApiVersionSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Api Version Set ID
func (id ApiVersionSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apiVersionSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.VersionSetId)
}

// Segments returns a slice of Resource ID Segments which comprise this Api Version Set ID
func (id ApiVersionSetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticApiVersionSets", "apiVersionSets", "apiVersionSets"),
		resourceids.UserSpecifiedSegment("versionSetId", "versionSetIdValue"),
	}
}

// String returns a human-readable description of this Api Version Set ID
func (id ApiVersionSetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Version Set: %q", id.VersionSetId),
	}
	return fmt.Sprintf("Api Version Set (%s)", strings.Join(components, "\n"))
}
