package httprouteconfig

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&HTTPRouteConfigId{})
}

var _ resourceids.ResourceId = &HTTPRouteConfigId{}

// HTTPRouteConfigId is a struct representing the Resource ID for a HTTP Route Config
type HTTPRouteConfigId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ManagedEnvironmentName string
	HttpRouteConfigName    string
}

// NewHTTPRouteConfigID returns a new HTTPRouteConfigId struct
func NewHTTPRouteConfigID(subscriptionId string, resourceGroupName string, managedEnvironmentName string, httpRouteConfigName string) HTTPRouteConfigId {
	return HTTPRouteConfigId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ManagedEnvironmentName: managedEnvironmentName,
		HttpRouteConfigName:    httpRouteConfigName,
	}
}

// ParseHTTPRouteConfigID parses 'input' into a HTTPRouteConfigId
func ParseHTTPRouteConfigID(input string) (*HTTPRouteConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HTTPRouteConfigId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HTTPRouteConfigId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHTTPRouteConfigIDInsensitively parses 'input' case-insensitively into a HTTPRouteConfigId
// note: this method should only be used for API response data and not user input
func ParseHTTPRouteConfigIDInsensitively(input string) (*HTTPRouteConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HTTPRouteConfigId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HTTPRouteConfigId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HTTPRouteConfigId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ManagedEnvironmentName, ok = input.Parsed["managedEnvironmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedEnvironmentName", input)
	}

	if id.HttpRouteConfigName, ok = input.Parsed["httpRouteConfigName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "httpRouteConfigName", input)
	}

	return nil
}

// ValidateHTTPRouteConfigID checks that 'input' can be parsed as a HTTP Route Config ID
func ValidateHTTPRouteConfigID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHTTPRouteConfigID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted HTTP Route Config ID
func (id HTTPRouteConfigId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s/httpRouteConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName, id.HttpRouteConfigName)
}

// Segments returns a slice of Resource ID Segments which comprise this HTTP Route Config ID
func (id HTTPRouteConfigId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticManagedEnvironments", "managedEnvironments", "managedEnvironments"),
		resourceids.UserSpecifiedSegment("managedEnvironmentName", "managedEnvironmentName"),
		resourceids.StaticSegment("staticHttpRouteConfigs", "httpRouteConfigs", "httpRouteConfigs"),
		resourceids.UserSpecifiedSegment("httpRouteConfigName", "httpRouteConfigName"),
	}
}

// String returns a human-readable description of this HTTP Route Config ID
func (id HTTPRouteConfigId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Environment Name: %q", id.ManagedEnvironmentName),
		fmt.Sprintf("Http Route Config Name: %q", id.HttpRouteConfigName),
	}
	return fmt.Sprintf("HTTP Route Config (%s)", strings.Join(components, "\n"))
}
