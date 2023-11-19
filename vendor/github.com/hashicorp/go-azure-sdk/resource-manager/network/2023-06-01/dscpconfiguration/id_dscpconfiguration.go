package dscpconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DscpConfigurationId{}

// DscpConfigurationId is a struct representing the Resource ID for a Dscp Configuration
type DscpConfigurationId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DscpConfigurationName string
}

// NewDscpConfigurationID returns a new DscpConfigurationId struct
func NewDscpConfigurationID(subscriptionId string, resourceGroupName string, dscpConfigurationName string) DscpConfigurationId {
	return DscpConfigurationId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		DscpConfigurationName: dscpConfigurationName,
	}
}

// ParseDscpConfigurationID parses 'input' into a DscpConfigurationId
func ParseDscpConfigurationID(input string) (*DscpConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(DscpConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DscpConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DscpConfigurationName, ok = parsed.Parsed["dscpConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dscpConfigurationName", *parsed)
	}

	return &id, nil
}

// ParseDscpConfigurationIDInsensitively parses 'input' case-insensitively into a DscpConfigurationId
// note: this method should only be used for API response data and not user input
func ParseDscpConfigurationIDInsensitively(input string) (*DscpConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(DscpConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DscpConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DscpConfigurationName, ok = parsed.Parsed["dscpConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dscpConfigurationName", *parsed)
	}

	return &id, nil
}

// ValidateDscpConfigurationID checks that 'input' can be parsed as a Dscp Configuration ID
func ValidateDscpConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDscpConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dscp Configuration ID
func (id DscpConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dscpConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DscpConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dscp Configuration ID
func (id DscpConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDscpConfigurations", "dscpConfigurations", "dscpConfigurations"),
		resourceids.UserSpecifiedSegment("dscpConfigurationName", "dscpConfigurationValue"),
	}
}

// String returns a human-readable description of this Dscp Configuration ID
func (id DscpConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dscp Configuration Name: %q", id.DscpConfigurationName),
	}
	return fmt.Sprintf("Dscp Configuration (%s)", strings.Join(components, "\n"))
}
