package apiissuecomment

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CommentId{})
}

var _ resourceids.ResourceId = &CommentId{}

// CommentId is a struct representing the Resource ID for a Comment
type CommentId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ApiId             string
	IssueId           string
	CommentId         string
}

// NewCommentID returns a new CommentId struct
func NewCommentID(subscriptionId string, resourceGroupName string, serviceName string, apiId string, issueId string, commentId string) CommentId {
	return CommentId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ApiId:             apiId,
		IssueId:           issueId,
		CommentId:         commentId,
	}
}

// ParseCommentID parses 'input' into a CommentId
func ParseCommentID(input string) (*CommentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CommentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CommentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCommentIDInsensitively parses 'input' case-insensitively into a CommentId
// note: this method should only be used for API response data and not user input
func ParseCommentIDInsensitively(input string) (*CommentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CommentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CommentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CommentId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.CommentId, ok = input.Parsed["commentId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "commentId", input)
	}

	return nil
}

// ValidateCommentID checks that 'input' can be parsed as a Comment ID
func ValidateCommentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCommentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Comment ID
func (id CommentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/issues/%s/comments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.IssueId, id.CommentId)
}

// Segments returns a slice of Resource ID Segments which comprise this Comment ID
func (id CommentId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticComments", "comments", "comments"),
		resourceids.UserSpecifiedSegment("commentId", "commentId"),
	}
}

// String returns a human-readable description of this Comment ID
func (id CommentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Api: %q", id.ApiId),
		fmt.Sprintf("Issue: %q", id.IssueId),
		fmt.Sprintf("Comment: %q", id.CommentId),
	}
	return fmt.Sprintf("Comment (%s)", strings.Join(components, "\n"))
}
