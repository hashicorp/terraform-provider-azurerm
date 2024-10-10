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
	recaser.RegisterResourceId(&AppSettingId{})
}

var _ resourceids.ResourceId = &AppSettingId{}

// AppSettingId is a struct representing the Resource ID for a App Setting
type AppSettingId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	AppSettingKey     string
}

// NewAppSettingID returns a new AppSettingId struct
func NewAppSettingID(subscriptionId string, resourceGroupName string, siteName string, appSettingKey string) AppSettingId {
	return AppSettingId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		AppSettingKey:     appSettingKey,
	}
}

// ParseAppSettingID parses 'input' into a AppSettingId
func ParseAppSettingID(input string) (*AppSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AppSettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AppSettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAppSettingIDInsensitively parses 'input' case-insensitively into a AppSettingId
// note: this method should only be used for API response data and not user input
func ParseAppSettingIDInsensitively(input string) (*AppSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AppSettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AppSettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AppSettingId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AppSettingKey, ok = input.Parsed["appSettingKey"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "appSettingKey", input)
	}

	return nil
}

// ValidateAppSettingID checks that 'input' can be parsed as a App Setting ID
func ValidateAppSettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAppSettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted App Setting ID
func (id AppSettingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/config/configReferences/appSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.AppSettingKey)
}

// Segments returns a slice of Resource ID Segments which comprise this App Setting ID
func (id AppSettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticConfig", "config", "config"),
		resourceids.StaticSegment("staticConfigReferences", "configReferences", "configReferences"),
		resourceids.StaticSegment("staticAppSettings", "appSettings", "appSettings"),
		resourceids.UserSpecifiedSegment("appSettingKey", "appSettingKey"),
	}
}

// String returns a human-readable description of this App Setting ID
func (id AppSettingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("App Setting Key: %q", id.AppSettingKey),
	}
	return fmt.Sprintf("App Setting (%s)", strings.Join(components, "\n"))
}
