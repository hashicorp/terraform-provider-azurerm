package apps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = IotAppId{}

// IotAppId is a struct representing the Resource ID for a Iot App
type IotAppId struct {
	SubscriptionId    string
	ResourceGroupName string
	IotAppName        string
}

// NewIotAppID returns a new IotAppId struct
func NewIotAppID(subscriptionId string, resourceGroupName string, iotAppName string) IotAppId {
	return IotAppId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		IotAppName:        iotAppName,
	}
}

// ParseIotAppID parses 'input' into a IotAppId
func ParseIotAppID(input string) (*IotAppId, error) {
	parser := resourceids.NewParserFromResourceIdType(IotAppId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IotAppId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.IotAppName, ok = parsed.Parsed["iotAppName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "iotAppName", *parsed)
	}

	return &id, nil
}

// ParseIotAppIDInsensitively parses 'input' case-insensitively into a IotAppId
// note: this method should only be used for API response data and not user input
func ParseIotAppIDInsensitively(input string) (*IotAppId, error) {
	parser := resourceids.NewParserFromResourceIdType(IotAppId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IotAppId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.IotAppName, ok = parsed.Parsed["iotAppName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "iotAppName", *parsed)
	}

	return &id, nil
}

// ValidateIotAppID checks that 'input' can be parsed as a Iot App ID
func ValidateIotAppID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIotAppID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Iot App ID
func (id IotAppId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.IoTCentral/iotApps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IotAppName)
}

// Segments returns a slice of Resource ID Segments which comprise this Iot App ID
func (id IotAppId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftIoTCentral", "Microsoft.IoTCentral", "Microsoft.IoTCentral"),
		resourceids.StaticSegment("staticIotApps", "iotApps", "iotApps"),
		resourceids.UserSpecifiedSegment("iotAppName", "iotAppValue"),
	}
}

// String returns a human-readable description of this Iot App ID
func (id IotAppId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Iot App Name: %q", id.IotAppName),
	}
	return fmt.Sprintf("Iot App (%s)", strings.Join(components, "\n"))
}
