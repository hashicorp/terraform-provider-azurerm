package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ConfigReferenceAppSettingId{}

// ConfigReferenceAppSettingId is a struct representing the Resource ID for a Config Reference App Setting
type ConfigReferenceAppSettingId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SlotName          string
	AppSettingKey     string
}

// NewConfigReferenceAppSettingID returns a new ConfigReferenceAppSettingId struct
func NewConfigReferenceAppSettingID(subscriptionId string, resourceGroupName string, siteName string, slotName string, appSettingKey string) ConfigReferenceAppSettingId {
	return ConfigReferenceAppSettingId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SlotName:          slotName,
		AppSettingKey:     appSettingKey,
	}
}

// ParseConfigReferenceAppSettingID parses 'input' into a ConfigReferenceAppSettingId
func ParseConfigReferenceAppSettingID(input string) (*ConfigReferenceAppSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigReferenceAppSettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigReferenceAppSettingId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConfigReferenceAppSettingIDInsensitively parses 'input' case-insensitively into a ConfigReferenceAppSettingId
// note: this method should only be used for API response data and not user input
func ParseConfigReferenceAppSettingIDInsensitively(input string) (*ConfigReferenceAppSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigReferenceAppSettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigReferenceAppSettingId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConfigReferenceAppSettingId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SlotName, ok = input.Parsed["slotName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "slotName", input)
	}

	if id.AppSettingKey, ok = input.Parsed["appSettingKey"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "appSettingKey", input)
	}

	return nil
}

// ValidateConfigReferenceAppSettingID checks that 'input' can be parsed as a Config Reference App Setting ID
func ValidateConfigReferenceAppSettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigReferenceAppSettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Config Reference App Setting ID
func (id ConfigReferenceAppSettingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/config/configReferences/appSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.AppSettingKey)
}

// Segments returns a slice of Resource ID Segments which comprise this Config Reference App Setting ID
func (id ConfigReferenceAppSettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteValue"),
		resourceids.StaticSegment("staticSlots", "slots", "slots"),
		resourceids.UserSpecifiedSegment("slotName", "slotValue"),
		resourceids.StaticSegment("staticConfig", "config", "config"),
		resourceids.StaticSegment("staticConfigReferences", "configReferences", "configReferences"),
		resourceids.StaticSegment("staticAppSettings", "appSettings", "appSettings"),
		resourceids.UserSpecifiedSegment("appSettingKey", "appSettingKeyValue"),
	}
}

// String returns a human-readable description of this Config Reference App Setting ID
func (id ConfigReferenceAppSettingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("App Setting Key: %q", id.AppSettingKey),
	}
	return fmt.Sprintf("Config Reference App Setting (%s)", strings.Join(components, "\n"))
}
