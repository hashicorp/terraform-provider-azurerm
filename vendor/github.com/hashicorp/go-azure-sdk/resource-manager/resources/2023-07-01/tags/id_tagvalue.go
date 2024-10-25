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
	recaser.RegisterResourceId(&TagValueId{})
}

var _ resourceids.ResourceId = &TagValueId{}

// TagValueId is a struct representing the Resource ID for a Tag Value
type TagValueId struct {
	SubscriptionId string
	TagName        string
	TagValueName   string
}

// NewTagValueID returns a new TagValueId struct
func NewTagValueID(subscriptionId string, tagName string, tagValueName string) TagValueId {
	return TagValueId{
		SubscriptionId: subscriptionId,
		TagName:        tagName,
		TagValueName:   tagValueName,
	}
}

// ParseTagValueID parses 'input' into a TagValueId
func ParseTagValueID(input string) (*TagValueId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TagValueId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TagValueId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTagValueIDInsensitively parses 'input' case-insensitively into a TagValueId
// note: this method should only be used for API response data and not user input
func ParseTagValueIDInsensitively(input string) (*TagValueId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TagValueId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TagValueId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TagValueId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.TagName, ok = input.Parsed["tagName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tagName", input)
	}

	if id.TagValueName, ok = input.Parsed["tagValueName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tagValueName", input)
	}

	return nil
}

// ValidateTagValueID checks that 'input' can be parsed as a Tag Value ID
func ValidateTagValueID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTagValueID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Tag Value ID
func (id TagValueId) ID() string {
	fmtString := "/subscriptions/%s/tagNames/%s/tagValues/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.TagName, id.TagValueName)
}

// Segments returns a slice of Resource ID Segments which comprise this Tag Value ID
func (id TagValueId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticTagNames", "tagNames", "tagNames"),
		resourceids.UserSpecifiedSegment("tagName", "tagName"),
		resourceids.StaticSegment("staticTagValues", "tagValues", "tagValues"),
		resourceids.UserSpecifiedSegment("tagValueName", "tagValueName"),
	}
}

// String returns a human-readable description of this Tag Value ID
func (id TagValueId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Tag Name: %q", id.TagName),
		fmt.Sprintf("Tag Value Name: %q", id.TagValueName),
	}
	return fmt.Sprintf("Tag Value (%s)", strings.Join(components, "\n"))
}
