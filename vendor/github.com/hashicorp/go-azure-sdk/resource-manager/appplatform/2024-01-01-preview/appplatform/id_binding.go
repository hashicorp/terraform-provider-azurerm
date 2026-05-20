package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BindingId{})
}

var _ resourceids.ResourceId = &BindingId{}

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
	parser := resourceids.NewParserFromResourceIdType(&BindingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BindingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBindingIDInsensitively parses 'input' case-insensitively into a BindingId
// note: this method should only be used for API response data and not user input
func ParseBindingIDInsensitively(input string) (*BindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BindingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BindingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BindingId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpringName, ok = input.Parsed["springName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "springName", input)
	}

	if id.AppName, ok = input.Parsed["appName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "appName", input)
	}

	if id.BindingName, ok = input.Parsed["bindingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "bindingName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticApps", "apps", "apps"),
		resourceids.UserSpecifiedSegment("appName", "appName"),
		resourceids.StaticSegment("staticBindings", "bindings", "bindings"),
		resourceids.UserSpecifiedSegment("bindingName", "bindingName"),
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
