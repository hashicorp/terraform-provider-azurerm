package apiissueattachment

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApiIssueId{})
}

var _ resourceids.ResourceId = &ApiIssueId{}

// ApiIssueId is a struct representing the Resource ID for a Api Issue
type ApiIssueId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ApiId             string
	IssueId           string
}

// NewApiIssueID returns a new ApiIssueId struct
func NewApiIssueID(subscriptionId string, resourceGroupName string, serviceName string, apiId string, issueId string) ApiIssueId {
	return ApiIssueId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ApiId:             apiId,
		IssueId:           issueId,
	}
}

// ParseApiIssueID parses 'input' into a ApiIssueId
func ParseApiIssueID(input string) (*ApiIssueId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiIssueId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiIssueId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApiIssueIDInsensitively parses 'input' case-insensitively into a ApiIssueId
// note: this method should only be used for API response data and not user input
func ParseApiIssueIDInsensitively(input string) (*ApiIssueId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiIssueId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiIssueId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApiIssueId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.IssueId, ok = input.Parsed["issueId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "issueId", input)
	}

	return nil
}

// ValidateApiIssueID checks that 'input' can be parsed as a Api Issue ID
func ValidateApiIssueID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApiIssueID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Api Issue ID
func (id ApiIssueId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/issues/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.IssueId)
}

// Segments returns a slice of Resource ID Segments which comprise this Api Issue ID
func (id ApiIssueId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticIssues", "issues", "issues"),
		resourceids.UserSpecifiedSegment("issueId", "issueId"),
	}
}

// String returns a human-readable description of this Api Issue ID
func (id ApiIssueId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Api: %q", id.ApiId),
		fmt.Sprintf("Issue: %q", id.IssueId),
	}
	return fmt.Sprintf("Api Issue (%s)", strings.Join(components, "\n"))
}
