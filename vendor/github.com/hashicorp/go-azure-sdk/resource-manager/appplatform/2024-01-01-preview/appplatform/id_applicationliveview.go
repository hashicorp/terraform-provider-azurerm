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
	recaser.RegisterResourceId(&ApplicationLiveViewId{})
}

var _ resourceids.ResourceId = &ApplicationLiveViewId{}

// ApplicationLiveViewId is a struct representing the Resource ID for a Application Live View
type ApplicationLiveViewId struct {
	SubscriptionId          string
	ResourceGroupName       string
	SpringName              string
	ApplicationLiveViewName string
}

// NewApplicationLiveViewID returns a new ApplicationLiveViewId struct
func NewApplicationLiveViewID(subscriptionId string, resourceGroupName string, springName string, applicationLiveViewName string) ApplicationLiveViewId {
	return ApplicationLiveViewId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		SpringName:              springName,
		ApplicationLiveViewName: applicationLiveViewName,
	}
}

// ParseApplicationLiveViewID parses 'input' into a ApplicationLiveViewId
func ParseApplicationLiveViewID(input string) (*ApplicationLiveViewId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationLiveViewId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationLiveViewId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApplicationLiveViewIDInsensitively parses 'input' case-insensitively into a ApplicationLiveViewId
// note: this method should only be used for API response data and not user input
func ParseApplicationLiveViewIDInsensitively(input string) (*ApplicationLiveViewId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationLiveViewId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationLiveViewId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApplicationLiveViewId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ApplicationLiveViewName, ok = input.Parsed["applicationLiveViewName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "applicationLiveViewName", input)
	}

	return nil
}

// ValidateApplicationLiveViewID checks that 'input' can be parsed as a Application Live View ID
func ValidateApplicationLiveViewID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationLiveViewID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Live View ID
func (id ApplicationLiveViewId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/applicationLiveViews/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ApplicationLiveViewName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Live View ID
func (id ApplicationLiveViewId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticApplicationLiveViews", "applicationLiveViews", "applicationLiveViews"),
		resourceids.UserSpecifiedSegment("applicationLiveViewName", "applicationLiveViewName"),
	}
}

// String returns a human-readable description of this Application Live View ID
func (id ApplicationLiveViewId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Application Live View Name: %q", id.ApplicationLiveViewName),
	}
	return fmt.Sprintf("Application Live View (%s)", strings.Join(components, "\n"))
}
