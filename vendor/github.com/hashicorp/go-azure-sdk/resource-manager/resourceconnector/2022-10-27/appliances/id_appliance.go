package appliances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ApplianceId{}

// ApplianceId is a struct representing the Resource ID for a Appliance
type ApplianceId struct {
	SubscriptionId    string
	ResourceGroupName string
	ApplianceName     string
}

// NewApplianceID returns a new ApplianceId struct
func NewApplianceID(subscriptionId string, resourceGroupName string, applianceName string) ApplianceId {
	return ApplianceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ApplianceName:     applianceName,
	}
}

// ParseApplianceID parses 'input' into a ApplianceId
func ParseApplianceID(input string) (*ApplianceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplianceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplianceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ApplianceName, ok = parsed.Parsed["applianceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "applianceName", *parsed)
	}

	return &id, nil
}

// ParseApplianceIDInsensitively parses 'input' case-insensitively into a ApplianceId
// note: this method should only be used for API response data and not user input
func ParseApplianceIDInsensitively(input string) (*ApplianceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplianceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplianceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ApplianceName, ok = parsed.Parsed["applianceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "applianceName", *parsed)
	}

	return &id, nil
}

// ValidateApplianceID checks that 'input' can be parsed as a Appliance ID
func ValidateApplianceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplianceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Appliance ID
func (id ApplianceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ResourceConnector/appliances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ApplianceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Appliance ID
func (id ApplianceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftResourceConnector", "Microsoft.ResourceConnector", "Microsoft.ResourceConnector"),
		resourceids.StaticSegment("staticAppliances", "appliances", "appliances"),
		resourceids.UserSpecifiedSegment("applianceName", "applianceValue"),
	}
}

// String returns a human-readable description of this Appliance ID
func (id ApplianceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Appliance Name: %q", id.ApplianceName),
	}
	return fmt.Sprintf("Appliance (%s)", strings.Join(components, "\n"))
}
