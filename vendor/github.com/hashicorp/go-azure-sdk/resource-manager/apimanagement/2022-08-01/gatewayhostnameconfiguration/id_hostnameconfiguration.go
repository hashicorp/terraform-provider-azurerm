package gatewayhostnameconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = HostnameConfigurationId{}

// HostnameConfigurationId is a struct representing the Resource ID for a Hostname Configuration
type HostnameConfigurationId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	GatewayId         string
	HcId              string
}

// NewHostnameConfigurationID returns a new HostnameConfigurationId struct
func NewHostnameConfigurationID(subscriptionId string, resourceGroupName string, serviceName string, gatewayId string, hcId string) HostnameConfigurationId {
	return HostnameConfigurationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		GatewayId:         gatewayId,
		HcId:              hcId,
	}
}

// ParseHostnameConfigurationID parses 'input' into a HostnameConfigurationId
func ParseHostnameConfigurationID(input string) (*HostnameConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(HostnameConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HostnameConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.GatewayId, ok = parsed.Parsed["gatewayId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "gatewayId", *parsed)
	}

	if id.HcId, ok = parsed.Parsed["hcId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hcId", *parsed)
	}

	return &id, nil
}

// ParseHostnameConfigurationIDInsensitively parses 'input' case-insensitively into a HostnameConfigurationId
// note: this method should only be used for API response data and not user input
func ParseHostnameConfigurationIDInsensitively(input string) (*HostnameConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(HostnameConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HostnameConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.GatewayId, ok = parsed.Parsed["gatewayId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "gatewayId", *parsed)
	}

	if id.HcId, ok = parsed.Parsed["hcId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hcId", *parsed)
	}

	return &id, nil
}

// ValidateHostnameConfigurationID checks that 'input' can be parsed as a Hostname Configuration ID
func ValidateHostnameConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHostnameConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hostname Configuration ID
func (id HostnameConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/gateways/%s/hostnameConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.GatewayId, id.HcId)
}

// Segments returns a slice of Resource ID Segments which comprise this Hostname Configuration ID
func (id HostnameConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticGateways", "gateways", "gateways"),
		resourceids.UserSpecifiedSegment("gatewayId", "gatewayIdValue"),
		resourceids.StaticSegment("staticHostnameConfigurations", "hostnameConfigurations", "hostnameConfigurations"),
		resourceids.UserSpecifiedSegment("hcId", "hcIdValue"),
	}
}

// String returns a human-readable description of this Hostname Configuration ID
func (id HostnameConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Gateway: %q", id.GatewayId),
		fmt.Sprintf("Hc: %q", id.HcId),
	}
	return fmt.Sprintf("Hostname Configuration (%s)", strings.Join(components, "\n"))
}
