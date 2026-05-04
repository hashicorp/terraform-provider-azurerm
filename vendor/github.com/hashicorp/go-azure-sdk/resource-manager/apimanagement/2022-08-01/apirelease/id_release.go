package apirelease

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ReleaseId{})
}

var _ resourceids.ResourceId = &ReleaseId{}

// ReleaseId is a struct representing the Resource ID for a Release
type ReleaseId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ApiId             string
	ReleaseId         string
}

// NewReleaseID returns a new ReleaseId struct
func NewReleaseID(subscriptionId string, resourceGroupName string, serviceName string, apiId string, releaseId string) ReleaseId {
	return ReleaseId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ApiId:             apiId,
		ReleaseId:         releaseId,
	}
}

// ParseReleaseID parses 'input' into a ReleaseId
func ParseReleaseID(input string) (*ReleaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReleaseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReleaseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseReleaseIDInsensitively parses 'input' case-insensitively into a ReleaseId
// note: this method should only be used for API response data and not user input
func ParseReleaseIDInsensitively(input string) (*ReleaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReleaseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReleaseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ReleaseId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ReleaseId, ok = input.Parsed["releaseId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "releaseId", input)
	}

	return nil
}

// ValidateReleaseID checks that 'input' can be parsed as a Release ID
func ValidateReleaseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReleaseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Release ID
func (id ReleaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/releases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.ReleaseId)
}

// Segments returns a slice of Resource ID Segments which comprise this Release ID
func (id ReleaseId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticReleases", "releases", "releases"),
		resourceids.UserSpecifiedSegment("releaseId", "releaseId"),
	}
}

// String returns a human-readable description of this Release ID
func (id ReleaseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Api: %q", id.ApiId),
		fmt.Sprintf("Release: %q", id.ReleaseId),
	}
	return fmt.Sprintf("Release (%s)", strings.Join(components, "\n"))
}
