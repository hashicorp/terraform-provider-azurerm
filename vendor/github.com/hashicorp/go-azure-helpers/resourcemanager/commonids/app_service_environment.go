// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &AppServiceEnvironmentId{}

// AppServiceEnvironmentId is a struct representing the Resource ID for a App Service Environment
type AppServiceEnvironmentId struct {
	SubscriptionId         string
	ResourceGroupName      string
	HostingEnvironmentName string
}

// NewAppServiceEnvironmentID returns a new AppServiceEnvironmentId struct
func NewAppServiceEnvironmentID(subscriptionId string, resourceGroupName string, hostingEnvironmentName string) AppServiceEnvironmentId {
	return AppServiceEnvironmentId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		HostingEnvironmentName: hostingEnvironmentName,
	}
}

// ParseAppServiceEnvironmentID parses 'input' into a AppServiceEnvironmentId
func ParseAppServiceEnvironmentID(input string) (*AppServiceEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AppServiceEnvironmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AppServiceEnvironmentId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAppServiceEnvironmentIDInsensitively parses 'input' case-insensitively into a AppServiceEnvironmentId
// note: this method should only be used for API response data and not user input
func ParseAppServiceEnvironmentIDInsensitively(input string) (*AppServiceEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AppServiceEnvironmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AppServiceEnvironmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ValidateAppServiceEnvironmentID checks that 'input' can be parsed as an App Service Environment ID
func ValidateAppServiceEnvironmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAppServiceEnvironmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

func (id *AppServiceEnvironmentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.HostingEnvironmentName, ok = input.Parsed["hostingEnvironmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hostingEnvironmentName", input)
	}

	return nil
}

// ID returns the formatted App Service Environment ID
func (id AppServiceEnvironmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/hostingEnvironments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostingEnvironmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this App Service Environment ID
func (id AppServiceEnvironmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticHostingEnvironments", "hostingEnvironments", "hostingEnvironments"),
		resourceids.UserSpecifiedSegment("hostingEnvironmentName", "hostingEnvironmentValue"),
	}
}

// String returns a human-readable description of this App Service Environment ID
func (id AppServiceEnvironmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("App Service Environment Name: %q", id.HostingEnvironmentName),
	}
	return fmt.Sprintf("App Service Environment (%s)", strings.Join(components, "\n"))
}
