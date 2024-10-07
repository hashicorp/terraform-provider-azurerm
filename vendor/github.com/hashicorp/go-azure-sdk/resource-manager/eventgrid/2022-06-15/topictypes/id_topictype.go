package topictypes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TopicTypeId{})
}

var _ resourceids.ResourceId = &TopicTypeId{}

// TopicTypeId is a struct representing the Resource ID for a Topic Type
type TopicTypeId struct {
	TopicTypeName string
}

// NewTopicTypeID returns a new TopicTypeId struct
func NewTopicTypeID(topicTypeName string) TopicTypeId {
	return TopicTypeId{
		TopicTypeName: topicTypeName,
	}
}

// ParseTopicTypeID parses 'input' into a TopicTypeId
func ParseTopicTypeID(input string) (*TopicTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TopicTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TopicTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTopicTypeIDInsensitively parses 'input' case-insensitively into a TopicTypeId
// note: this method should only be used for API response data and not user input
func ParseTopicTypeIDInsensitively(input string) (*TopicTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TopicTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TopicTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TopicTypeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.TopicTypeName, ok = input.Parsed["topicTypeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "topicTypeName", input)
	}

	return nil
}

// ValidateTopicTypeID checks that 'input' can be parsed as a Topic Type ID
func ValidateTopicTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTopicTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Topic Type ID
func (id TopicTypeId) ID() string {
	fmtString := "/providers/Microsoft.EventGrid/topicTypes/%s"
	return fmt.Sprintf(fmtString, id.TopicTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Topic Type ID
func (id TopicTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticTopicTypes", "topicTypes", "topicTypes"),
		resourceids.UserSpecifiedSegment("topicTypeName", "topicTypeName"),
	}
}

// String returns a human-readable description of this Topic Type ID
func (id TopicTypeId) String() string {
	components := []string{
		fmt.Sprintf("Topic Type Name: %q", id.TopicTypeName),
	}
	return fmt.Sprintf("Topic Type (%s)", strings.Join(components, "\n"))
}
