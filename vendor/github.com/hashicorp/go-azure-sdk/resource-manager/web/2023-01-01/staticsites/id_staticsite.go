package staticsites

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StaticSiteId{})
}

var _ resourceids.ResourceId = &StaticSiteId{}

// StaticSiteId is a struct representing the Resource ID for a Static Site
type StaticSiteId struct {
	SubscriptionId    string
	ResourceGroupName string
	StaticSiteName    string
}

// NewStaticSiteID returns a new StaticSiteId struct
func NewStaticSiteID(subscriptionId string, resourceGroupName string, staticSiteName string) StaticSiteId {
	return StaticSiteId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StaticSiteName:    staticSiteName,
	}
}

// ParseStaticSiteID parses 'input' into a StaticSiteId
func ParseStaticSiteID(input string) (*StaticSiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StaticSiteId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StaticSiteId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStaticSiteIDInsensitively parses 'input' case-insensitively into a StaticSiteId
// note: this method should only be used for API response data and not user input
func ParseStaticSiteIDInsensitively(input string) (*StaticSiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StaticSiteId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StaticSiteId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StaticSiteId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StaticSiteName, ok = input.Parsed["staticSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "staticSiteName", input)
	}

	return nil
}

// ValidateStaticSiteID checks that 'input' can be parsed as a Static Site ID
func ValidateStaticSiteID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStaticSiteID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Static Site ID
func (id StaticSiteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/staticSites/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName)
}

// Segments returns a slice of Resource ID Segments which comprise this Static Site ID
func (id StaticSiteId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticStaticSites", "staticSites", "staticSites"),
		resourceids.UserSpecifiedSegment("staticSiteName", "staticSiteName"),
	}
}

// String returns a human-readable description of this Static Site ID
func (id StaticSiteId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Static Site Name: %q", id.StaticSiteName),
	}
	return fmt.Sprintf("Static Site (%s)", strings.Join(components, "\n"))
}
