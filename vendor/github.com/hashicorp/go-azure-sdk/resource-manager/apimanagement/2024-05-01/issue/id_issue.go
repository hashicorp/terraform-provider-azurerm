package issue

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&IssueId{})
}

var _ resourceids.ResourceId = &IssueId{}

// IssueId is a struct representing the Resource ID for a Issue
type IssueId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	IssueId           string
}

// NewIssueID returns a new IssueId struct
func NewIssueID(subscriptionId string, resourceGroupName string, serviceName string, issueId string) IssueId {
	return IssueId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		IssueId:           issueId,
	}
}

// ParseIssueID parses 'input' into a IssueId
func ParseIssueID(input string) (*IssueId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IssueId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IssueId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseIssueIDInsensitively parses 'input' case-insensitively into a IssueId
// note: this method should only be used for API response data and not user input
func ParseIssueIDInsensitively(input string) (*IssueId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IssueId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IssueId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *IssueId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.IssueId, ok = input.Parsed["issueId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "issueId", input)
	}

	return nil
}

// ValidateIssueID checks that 'input' can be parsed as a Issue ID
func ValidateIssueID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIssueID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Issue ID
func (id IssueId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/issues/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.IssueId)
}

// Segments returns a slice of Resource ID Segments which comprise this Issue ID
func (id IssueId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticIssues", "issues", "issues"),
		resourceids.UserSpecifiedSegment("issueId", "issueId"),
	}
}

// String returns a human-readable description of this Issue ID
func (id IssueId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Issue: %q", id.IssueId),
	}
	return fmt.Sprintf("Issue (%s)", strings.Join(components, "\n"))
}
