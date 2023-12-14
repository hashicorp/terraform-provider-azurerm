// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = AppServiceId{}

// AppServiceId is a struct representing the Resource ID for an App Service
type AppServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
}

// NewAppServiceID returns a new AppServiceId struct
func NewAppServiceID(subscriptionId string, resourceGroupName string, siteName string) AppServiceId {
	return AppServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
	}
}

// ParseAppServiceID parses 'input' into a AppServiceId
func ParseAppServiceID(input string) (*AppServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(AppServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AppServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SiteName, ok = parsed.Parsed["siteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "siteName", *parsed)
	}

	return &id, nil
}

// ParseAppServiceIDInsensitively parses 'input' case-insensitively into a AppServiceId
// note: this method should only be used for API response data and not user input
func ParseAppServiceIDInsensitively(input string) (*AppServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(AppServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AppServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SiteName, ok = parsed.Parsed["siteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "siteName", *parsed)
	}

	return &id, nil
}

// ValidateAppServiceID checks that 'input' can be parsed as a App Service ID
func ValidateAppServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAppServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted App Service ID
func (id AppServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName)
}

// Segments returns a slice of Resource ID Segments which comprise this App Service ID
func (id AppServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteValue"),
	}
}

// String returns a human-readable description of this App Service ID
func (id AppServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
	}
	return fmt.Sprintf("App Service (%s)", strings.Join(components, "\n"))
}
