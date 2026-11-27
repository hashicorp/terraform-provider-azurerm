package apiversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApiId{})
}

var _ resourceids.ResourceId = &ApiId{}

// ApiId is a struct representing the Resource ID for a Api
type ApiId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceName     string
	ApiName           string
}

// NewApiID returns a new ApiId struct
func NewApiID(subscriptionId string, resourceGroupName string, serviceName string, workspaceName string, apiName string) ApiId {
	return ApiId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceName:     workspaceName,
		ApiName:           apiName,
	}
}

// ParseApiID parses 'input' into a ApiId
func ParseApiID(input string) (*ApiId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApiIDInsensitively parses 'input' case-insensitively into a ApiId
// note: this method should only be used for API response data and not user input
func ParseApiIDInsensitively(input string) (*ApiId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApiId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.WorkspaceName, ok = input.Parsed["workspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", input)
	}

	if id.ApiName, ok = input.Parsed["apiName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiName", input)
	}

	return nil
}

// ValidateApiID checks that 'input' can be parsed as a Api ID
func ValidateApiID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApiID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Api ID
func (id ApiId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiCenter/services/%s/workspaces/%s/apis/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceName, id.ApiName)
}

// Segments returns a slice of Resource ID Segments which comprise this Api ID
func (id ApiId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiCenter", "Microsoft.ApiCenter", "Microsoft.ApiCenter"),
		resourceids.StaticSegment("staticServices", "services", "services"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
		resourceids.StaticSegment("staticApis", "apis", "apis"),
		resourceids.UserSpecifiedSegment("apiName", "apiName"),
	}
}

// String returns a human-readable description of this Api ID
func (id ApiId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Api Name: %q", id.ApiName),
	}
	return fmt.Sprintf("Api (%s)", strings.Join(components, "\n"))
}
