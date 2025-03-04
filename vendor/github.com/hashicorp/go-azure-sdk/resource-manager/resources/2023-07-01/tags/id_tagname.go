package tags

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TagNameId{})
}

var _ resourceids.ResourceId = &TagNameId{}

// TagNameId is a struct representing the Resource ID for a Tag Name
type TagNameId struct {
	SubscriptionId string
	TagName        string
}

// NewTagNameID returns a new TagNameId struct
func NewTagNameID(subscriptionId string, tagName string) TagNameId {
	return TagNameId{
		SubscriptionId: subscriptionId,
		TagName:        tagName,
	}
}

// ParseTagNameID parses 'input' into a TagNameId
func ParseTagNameID(input string) (*TagNameId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TagNameId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TagNameId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTagNameIDInsensitively parses 'input' case-insensitively into a TagNameId
// note: this method should only be used for API response data and not user input
func ParseTagNameIDInsensitively(input string) (*TagNameId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TagNameId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TagNameId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TagNameId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.TagName, ok = input.Parsed["tagName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tagName", input)
	}

	return nil
}

// ValidateTagNameID checks that 'input' can be parsed as a Tag Name ID
func ValidateTagNameID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTagNameID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Tag Name ID
func (id TagNameId) ID() string {
	fmtString := "/subscriptions/%s/tagNames/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.TagName)
}

// Segments returns a slice of Resource ID Segments which comprise this Tag Name ID
func (id TagNameId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticTagNames", "tagNames", "tagNames"),
		resourceids.UserSpecifiedSegment("tagName", "tagName"),
	}
}

// String returns a human-readable description of this Tag Name ID
func (id TagNameId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Tag Name: %q", id.TagName),
	}
	return fmt.Sprintf("Tag Name (%s)", strings.Join(components, "\n"))
}
