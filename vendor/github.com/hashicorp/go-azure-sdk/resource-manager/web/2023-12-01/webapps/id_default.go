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
	recaser.RegisterResourceId(&DefaultId{})
}

var _ resourceids.ResourceId = &DefaultId{}

// DefaultId is a struct representing the Resource ID for a Default
type DefaultId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	DefaultName       string
	KeyName           string
}

// NewDefaultID returns a new DefaultId struct
func NewDefaultID(subscriptionId string, resourceGroupName string, siteName string, defaultName string, keyName string) DefaultId {
	return DefaultId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		DefaultName:       defaultName,
		KeyName:           keyName,
	}
}

// ParseDefaultID parses 'input' into a DefaultId
func ParseDefaultID(input string) (*DefaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DefaultId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DefaultId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDefaultIDInsensitively parses 'input' case-insensitively into a DefaultId
// note: this method should only be used for API response data and not user input
func ParseDefaultIDInsensitively(input string) (*DefaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DefaultId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DefaultId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DefaultId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DefaultName, ok = input.Parsed["defaultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "defaultName", input)
	}

	if id.KeyName, ok = input.Parsed["keyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "keyName", input)
	}

	return nil
}

// ValidateDefaultID checks that 'input' can be parsed as a Default ID
func ValidateDefaultID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDefaultID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Default ID
func (id DefaultId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/host/default/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.DefaultName, id.KeyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Default ID
func (id DefaultId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticHost", "host", "host"),
		resourceids.StaticSegment("staticDefault", "default", "default"),
		resourceids.UserSpecifiedSegment("defaultName", "defaultName"),
		resourceids.UserSpecifiedSegment("keyName", "keyName"),
	}
}

// String returns a human-readable description of this Default ID
func (id DefaultId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Default Name: %q", id.DefaultName),
		fmt.Sprintf("Key Name: %q", id.KeyName),
	}
	return fmt.Sprintf("Default (%s)", strings.Join(components, "\n"))
}
