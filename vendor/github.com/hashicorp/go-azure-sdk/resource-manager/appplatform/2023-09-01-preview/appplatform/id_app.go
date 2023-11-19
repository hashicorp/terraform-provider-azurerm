package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AppId{}

// AppId is a struct representing the Resource ID for a App
type AppId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	AppName           string
}

// NewAppID returns a new AppId struct
func NewAppID(subscriptionId string, resourceGroupName string, springName string, appName string) AppId {
	return AppId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		AppName:           appName,
	}
}

// ParseAppID parses 'input' into a AppId
func ParseAppID(input string) (*AppId, error) {
	parser := resourceids.NewParserFromResourceIdType(AppId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AppId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.AppName, ok = parsed.Parsed["appName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "appName", *parsed)
	}

	return &id, nil
}

// ParseAppIDInsensitively parses 'input' case-insensitively into a AppId
// note: this method should only be used for API response data and not user input
func ParseAppIDInsensitively(input string) (*AppId, error) {
	parser := resourceids.NewParserFromResourceIdType(AppId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AppId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.AppName, ok = parsed.Parsed["appName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "appName", *parsed)
	}

	return &id, nil
}

// ValidateAppID checks that 'input' can be parsed as a App ID
func ValidateAppID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAppID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted App ID
func (id AppId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.AppName)
}

// Segments returns a slice of Resource ID Segments which comprise this App ID
func (id AppId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springValue"),
		resourceids.StaticSegment("staticApps", "apps", "apps"),
		resourceids.UserSpecifiedSegment("appName", "appValue"),
	}
}

// String returns a human-readable description of this App ID
func (id AppId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("App Name: %q", id.AppName),
	}
	return fmt.Sprintf("App (%s)", strings.Join(components, "\n"))
}
