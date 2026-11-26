package firstpartyservicetags

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FirstPartyServiceTagId{})
}

var _ resourceids.ResourceId = &FirstPartyServiceTagId{}

// FirstPartyServiceTagId is a struct representing the Resource ID for a First Party Service Tag
type FirstPartyServiceTagId struct {
	SubscriptionId           string
	ResourceGroupName        string
	FirstPartyServiceTagName string
}

// NewFirstPartyServiceTagID returns a new FirstPartyServiceTagId struct
func NewFirstPartyServiceTagID(subscriptionId string, resourceGroupName string, firstPartyServiceTagName string) FirstPartyServiceTagId {
	return FirstPartyServiceTagId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		FirstPartyServiceTagName: firstPartyServiceTagName,
	}
}

// ParseFirstPartyServiceTagID parses 'input' into a FirstPartyServiceTagId
func ParseFirstPartyServiceTagID(input string) (*FirstPartyServiceTagId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FirstPartyServiceTagId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FirstPartyServiceTagId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFirstPartyServiceTagIDInsensitively parses 'input' case-insensitively into a FirstPartyServiceTagId
// note: this method should only be used for API response data and not user input
func ParseFirstPartyServiceTagIDInsensitively(input string) (*FirstPartyServiceTagId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FirstPartyServiceTagId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FirstPartyServiceTagId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FirstPartyServiceTagId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FirstPartyServiceTagName, ok = input.Parsed["firstPartyServiceTagName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "firstPartyServiceTagName", input)
	}

	return nil
}

// ValidateFirstPartyServiceTagID checks that 'input' can be parsed as a First Party Service Tag ID
func ValidateFirstPartyServiceTagID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFirstPartyServiceTagID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted First Party Service Tag ID
func (id FirstPartyServiceTagId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/firstPartyServiceTags/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FirstPartyServiceTagName)
}

// Segments returns a slice of Resource ID Segments which comprise this First Party Service Tag ID
func (id FirstPartyServiceTagId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticFirstPartyServiceTags", "firstPartyServiceTags", "firstPartyServiceTags"),
		resourceids.UserSpecifiedSegment("firstPartyServiceTagName", "firstPartyServiceTagName"),
	}
}

// String returns a human-readable description of this First Party Service Tag ID
func (id FirstPartyServiceTagId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("First Party Service Tag Name: %q", id.FirstPartyServiceTagName),
	}
	return fmt.Sprintf("First Party Service Tag (%s)", strings.Join(components, "\n"))
}
