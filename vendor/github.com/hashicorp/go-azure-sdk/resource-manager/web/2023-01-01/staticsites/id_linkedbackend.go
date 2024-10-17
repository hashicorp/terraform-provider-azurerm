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
	recaser.RegisterResourceId(&LinkedBackendId{})
}

var _ resourceids.ResourceId = &LinkedBackendId{}

// LinkedBackendId is a struct representing the Resource ID for a Linked Backend
type LinkedBackendId struct {
	SubscriptionId    string
	ResourceGroupName string
	StaticSiteName    string
	LinkedBackendName string
}

// NewLinkedBackendID returns a new LinkedBackendId struct
func NewLinkedBackendID(subscriptionId string, resourceGroupName string, staticSiteName string, linkedBackendName string) LinkedBackendId {
	return LinkedBackendId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StaticSiteName:    staticSiteName,
		LinkedBackendName: linkedBackendName,
	}
}

// ParseLinkedBackendID parses 'input' into a LinkedBackendId
func ParseLinkedBackendID(input string) (*LinkedBackendId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LinkedBackendId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LinkedBackendId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLinkedBackendIDInsensitively parses 'input' case-insensitively into a LinkedBackendId
// note: this method should only be used for API response data and not user input
func ParseLinkedBackendIDInsensitively(input string) (*LinkedBackendId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LinkedBackendId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LinkedBackendId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LinkedBackendId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.LinkedBackendName, ok = input.Parsed["linkedBackendName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "linkedBackendName", input)
	}

	return nil
}

// ValidateLinkedBackendID checks that 'input' can be parsed as a Linked Backend ID
func ValidateLinkedBackendID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLinkedBackendID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Linked Backend ID
func (id LinkedBackendId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/staticSites/%s/linkedBackends/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName, id.LinkedBackendName)
}

// Segments returns a slice of Resource ID Segments which comprise this Linked Backend ID
func (id LinkedBackendId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticStaticSites", "staticSites", "staticSites"),
		resourceids.UserSpecifiedSegment("staticSiteName", "staticSiteName"),
		resourceids.StaticSegment("staticLinkedBackends", "linkedBackends", "linkedBackends"),
		resourceids.UserSpecifiedSegment("linkedBackendName", "linkedBackendName"),
	}
}

// String returns a human-readable description of this Linked Backend ID
func (id LinkedBackendId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Static Site Name: %q", id.StaticSiteName),
		fmt.Sprintf("Linked Backend Name: %q", id.LinkedBackendName),
	}
	return fmt.Sprintf("Linked Backend (%s)", strings.Join(components, "\n"))
}
