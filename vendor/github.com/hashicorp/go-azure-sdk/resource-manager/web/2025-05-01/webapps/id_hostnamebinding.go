package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&HostNameBindingId{})
}

var _ resourceids.ResourceId = &HostNameBindingId{}

// HostNameBindingId is a struct representing the Resource ID for a Host Name Binding
type HostNameBindingId struct {
	SubscriptionId      string
	ResourceGroupName   string
	SiteName            string
	HostNameBindingName string
}

// NewHostNameBindingID returns a new HostNameBindingId struct
func NewHostNameBindingID(subscriptionId string, resourceGroupName string, siteName string, hostNameBindingName string) HostNameBindingId {
	return HostNameBindingId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		SiteName:            siteName,
		HostNameBindingName: hostNameBindingName,
	}
}

// ParseHostNameBindingID parses 'input' into a HostNameBindingId
func ParseHostNameBindingID(input string) (*HostNameBindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HostNameBindingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HostNameBindingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHostNameBindingIDInsensitively parses 'input' case-insensitively into a HostNameBindingId
// note: this method should only be used for API response data and not user input
func ParseHostNameBindingIDInsensitively(input string) (*HostNameBindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HostNameBindingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HostNameBindingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HostNameBindingId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.HostNameBindingName, ok = input.Parsed["hostNameBindingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hostNameBindingName", input)
	}

	return nil
}

// ValidateHostNameBindingID checks that 'input' can be parsed as a Host Name Binding ID
func ValidateHostNameBindingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHostNameBindingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Host Name Binding ID
func (id HostNameBindingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/hostNameBindings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.HostNameBindingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Host Name Binding ID
func (id HostNameBindingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticHostNameBindings", "hostNameBindings", "hostNameBindings"),
		resourceids.UserSpecifiedSegment("hostNameBindingName", "hostNameBindingName"),
	}
}

// String returns a human-readable description of this Host Name Binding ID
func (id HostNameBindingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Host Name Binding Name: %q", id.HostNameBindingName),
	}
	return fmt.Sprintf("Host Name Binding (%s)", strings.Join(components, "\n"))
}
