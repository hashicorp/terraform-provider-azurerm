package tagrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TagRuleId{})
}

var _ resourceids.ResourceId = &TagRuleId{}

// TagRuleId is a struct representing the Resource ID for a Tag Rule
type TagRuleId struct {
	SubscriptionId    string
	ResourceGroupName string
	MonitorName       string
	TagRuleName       string
}

// NewTagRuleID returns a new TagRuleId struct
func NewTagRuleID(subscriptionId string, resourceGroupName string, monitorName string, tagRuleName string) TagRuleId {
	return TagRuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MonitorName:       monitorName,
		TagRuleName:       tagRuleName,
	}
}

// ParseTagRuleID parses 'input' into a TagRuleId
func ParseTagRuleID(input string) (*TagRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TagRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TagRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTagRuleIDInsensitively parses 'input' case-insensitively into a TagRuleId
// note: this method should only be used for API response data and not user input
func ParseTagRuleIDInsensitively(input string) (*TagRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TagRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TagRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TagRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MonitorName, ok = input.Parsed["monitorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "monitorName", input)
	}

	if id.TagRuleName, ok = input.Parsed["tagRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tagRuleName", input)
	}

	return nil
}

// ValidateTagRuleID checks that 'input' can be parsed as a Tag Rule ID
func ValidateTagRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTagRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Tag Rule ID
func (id TagRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/NewRelic.Observability/monitors/%s/tagRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MonitorName, id.TagRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Tag Rule ID
func (id TagRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticNewRelicObservability", "NewRelic.Observability", "NewRelic.Observability"),
		resourceids.StaticSegment("staticMonitors", "monitors", "monitors"),
		resourceids.UserSpecifiedSegment("monitorName", "monitorName"),
		resourceids.StaticSegment("staticTagRules", "tagRules", "tagRules"),
		resourceids.UserSpecifiedSegment("tagRuleName", "tagRuleName"),
	}
}

// String returns a human-readable description of this Tag Rule ID
func (id TagRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Monitor Name: %q", id.MonitorName),
		fmt.Sprintf("Tag Rule Name: %q", id.TagRuleName),
	}
	return fmt.Sprintf("Tag Rule (%s)", strings.Join(components, "\n"))
}
