package portalconfig

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PortalConfigId{})
}

var _ resourceids.ResourceId = &PortalConfigId{}

// PortalConfigId is a struct representing the Resource ID for a Portal Config
type PortalConfigId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	PortalConfigId    string
}

// NewPortalConfigID returns a new PortalConfigId struct
func NewPortalConfigID(subscriptionId string, resourceGroupName string, serviceName string, portalConfigId string) PortalConfigId {
	return PortalConfigId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		PortalConfigId:    portalConfigId,
	}
}

// ParsePortalConfigID parses 'input' into a PortalConfigId
func ParsePortalConfigID(input string) (*PortalConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PortalConfigId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PortalConfigId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePortalConfigIDInsensitively parses 'input' case-insensitively into a PortalConfigId
// note: this method should only be used for API response data and not user input
func ParsePortalConfigIDInsensitively(input string) (*PortalConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PortalConfigId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PortalConfigId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PortalConfigId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.PortalConfigId, ok = input.Parsed["portalConfigId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "portalConfigId", input)
	}

	return nil
}

// ValidatePortalConfigID checks that 'input' can be parsed as a Portal Config ID
func ValidatePortalConfigID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePortalConfigID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Portal Config ID
func (id PortalConfigId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/portalConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.PortalConfigId)
}

// Segments returns a slice of Resource ID Segments which comprise this Portal Config ID
func (id PortalConfigId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticPortalConfigs", "portalConfigs", "portalConfigs"),
		resourceids.UserSpecifiedSegment("portalConfigId", "portalConfigId"),
	}
}

// String returns a human-readable description of this Portal Config ID
func (id PortalConfigId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Portal Config: %q", id.PortalConfigId),
	}
	return fmt.Sprintf("Portal Config (%s)", strings.Join(components, "\n"))
}
