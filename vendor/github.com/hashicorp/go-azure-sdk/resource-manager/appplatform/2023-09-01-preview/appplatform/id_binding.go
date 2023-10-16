package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = BindingId{}

// BindingId is a struct representing the Resource ID for a Binding
type BindingId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	AppName           string
	BindingName       string
}

// NewBindingID returns a new BindingId struct
func NewBindingID(subscriptionId string, resourceGroupName string, springName string, appName string, bindingName string) BindingId {
	return BindingId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		AppName:           appName,
		BindingName:       bindingName,
	}
}

// ParseBindingID parses 'input' into a BindingId
func ParseBindingID(input string) (*BindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(BindingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BindingId{}

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

	if id.BindingName, ok = parsed.Parsed["bindingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "bindingName", *parsed)
	}

	return &id, nil
}

// ParseBindingIDInsensitively parses 'input' case-insensitively into a BindingId
// note: this method should only be used for API response data and not user input
func ParseBindingIDInsensitively(input string) (*BindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(BindingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BindingId{}

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

	if id.BindingName, ok = parsed.Parsed["bindingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "bindingName", *parsed)
	}

	return &id, nil
}

// ValidateBindingID checks that 'input' can be parsed as a Binding ID
func ValidateBindingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBindingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Binding ID
func (id BindingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apps/%s/bindings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.AppName, id.BindingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Binding ID
func (id BindingId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticBindings", "bindings", "bindings"),
		resourceids.UserSpecifiedSegment("bindingName", "bindingValue"),
	}
}

// String returns a human-readable description of this Binding ID
func (id BindingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("App Name: %q", id.AppName),
		fmt.Sprintf("Binding Name: %q", id.BindingName),
	}
	return fmt.Sprintf("Binding (%s)", strings.Join(components, "\n"))
}
