package apiversionsets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApiVersionSetId{})
}

var _ resourceids.ResourceId = &ApiVersionSetId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ApiVersionSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiVersionSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApiVersionSetIDInsensitively parses 'input' case-insensitively into a ApiVersionSetId
// note: this method should only be used for API response data and not user input
func ParseApiVersionSetIDInsensitively(input string) (*ApiVersionSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiVersionSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiVersionSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApiVersionSetId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.VersionSetId, ok = input.Parsed["versionSetId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "versionSetId", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticApiVersionSets", "apiVersionSets", "apiVersionSets"),
		resourceids.UserSpecifiedSegment("versionSetId", "versionSetId"),
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
