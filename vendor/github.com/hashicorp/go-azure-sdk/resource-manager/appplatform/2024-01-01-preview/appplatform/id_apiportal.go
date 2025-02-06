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
	recaser.RegisterResourceId(&ApiPortalId{})
}

var _ resourceids.ResourceId = &ApiPortalId{}

// ApiPortalId is a struct representing the Resource ID for a Api Portal
type ApiPortalId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	ApiPortalName     string
}

// NewApiPortalID returns a new ApiPortalId struct
func NewApiPortalID(subscriptionId string, resourceGroupName string, springName string, apiPortalName string) ApiPortalId {
	return ApiPortalId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		ApiPortalName:     apiPortalName,
	}
}

// ParseApiPortalID parses 'input' into a ApiPortalId
func ParseApiPortalID(input string) (*ApiPortalId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiPortalId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiPortalId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApiPortalIDInsensitively parses 'input' case-insensitively into a ApiPortalId
// note: this method should only be used for API response data and not user input
func ParseApiPortalIDInsensitively(input string) (*ApiPortalId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiPortalId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiPortalId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApiPortalId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ApiPortalName, ok = input.Parsed["apiPortalName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiPortalName", input)
	}

	return nil
}

// ValidateApiPortalID checks that 'input' can be parsed as a Api Portal ID
func ValidateApiPortalID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApiPortalID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Api Portal ID
func (id ApiPortalId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apiPortals/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ApiPortalName)
}

// Segments returns a slice of Resource ID Segments which comprise this Api Portal ID
func (id ApiPortalId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticApiPortals", "apiPortals", "apiPortals"),
		resourceids.UserSpecifiedSegment("apiPortalName", "apiPortalName"),
	}
}

// String returns a human-readable description of this Api Portal ID
func (id ApiPortalId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Api Portal Name: %q", id.ApiPortalName),
	}
	return fmt.Sprintf("Api Portal (%s)", strings.Join(components, "\n"))
}
