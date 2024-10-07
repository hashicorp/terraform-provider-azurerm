package partnertopics

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PartnerTopicId{})
}

var _ resourceids.ResourceId = &PartnerTopicId{}

// PartnerTopicId is a struct representing the Resource ID for a Partner Topic
type PartnerTopicId struct {
	SubscriptionId    string
	ResourceGroupName string
	PartnerTopicName  string
}

// NewPartnerTopicID returns a new PartnerTopicId struct
func NewPartnerTopicID(subscriptionId string, resourceGroupName string, partnerTopicName string) PartnerTopicId {
	return PartnerTopicId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		PartnerTopicName:  partnerTopicName,
	}
}

// ParsePartnerTopicID parses 'input' into a PartnerTopicId
func ParsePartnerTopicID(input string) (*PartnerTopicId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PartnerTopicId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PartnerTopicId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePartnerTopicIDInsensitively parses 'input' case-insensitively into a PartnerTopicId
// note: this method should only be used for API response data and not user input
func ParsePartnerTopicIDInsensitively(input string) (*PartnerTopicId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PartnerTopicId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PartnerTopicId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PartnerTopicId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PartnerTopicName, ok = input.Parsed["partnerTopicName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "partnerTopicName", input)
	}

	return nil
}

// ValidatePartnerTopicID checks that 'input' can be parsed as a Partner Topic ID
func ValidatePartnerTopicID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePartnerTopicID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Partner Topic ID
func (id PartnerTopicId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/partnerTopics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PartnerTopicName)
}

// Segments returns a slice of Resource ID Segments which comprise this Partner Topic ID
func (id PartnerTopicId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticPartnerTopics", "partnerTopics", "partnerTopics"),
		resourceids.UserSpecifiedSegment("partnerTopicName", "partnerTopicName"),
	}
}

// String returns a human-readable description of this Partner Topic ID
func (id PartnerTopicId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Partner Topic Name: %q", id.PartnerTopicName),
	}
	return fmt.Sprintf("Partner Topic (%s)", strings.Join(components, "\n"))
}
