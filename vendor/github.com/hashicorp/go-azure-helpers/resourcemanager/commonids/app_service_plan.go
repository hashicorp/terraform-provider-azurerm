// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &AppServicePlanId{}

// AppServicePlanId is a struct representing the Resource ID for an App Service Plan
type AppServicePlanId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerFarmName    string
}

// NewAppServicePlanID returns a new AppServicePlanId struct
func NewAppServicePlanID(subscriptionId string, resourceGroupName string, serverFarmName string) AppServicePlanId {
	return AppServicePlanId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerFarmName:    serverFarmName,
	}
}

// ParseAppServicePlanID parses 'input' into a AppServicePlanId
func ParseAppServicePlanID(input string) (*AppServicePlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AppServicePlanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AppServicePlanId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAppServicePlanIDInsensitively parses 'input' case-insensitively into a AppServicePlanId
// note: this method should only be used for API response data and not user input
func ParseAppServicePlanIDInsensitively(input string) (*AppServicePlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AppServicePlanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AppServicePlanId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AppServicePlanId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServerFarmName, ok = input.Parsed["serverFarmName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serverFarmName", input)
	}

	return nil
}

// ValidateAppServicePlanID checks that 'input' can be parsed as an App Service Plan ID
func ValidateAppServicePlanID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAppServicePlanID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted App Service Plan ID
func (id AppServicePlanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/serverFarms/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerFarmName)
}

// Segments returns a slice of Resource ID Segments which comprise this App Service Plan ID
func (id AppServicePlanId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticServerFarms", "serverFarms", "serverFarms"),
		resourceids.UserSpecifiedSegment("serverFarmName", "serverFarmValue"),
	}
}

// String returns a human-readable description of this App Service Plan ID
func (id AppServicePlanId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Farm Name: %q", id.ServerFarmName),
	}
	return fmt.Sprintf("App Service Plan (%s)", strings.Join(components, "\n"))
}
