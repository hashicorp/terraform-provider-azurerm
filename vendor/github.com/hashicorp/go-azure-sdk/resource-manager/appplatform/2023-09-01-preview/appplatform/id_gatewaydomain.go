package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = GatewayDomainId{}

// GatewayDomainId is a struct representing the Resource ID for a Gateway Domain
type GatewayDomainId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	GatewayName       string
	DomainName        string
}

// NewGatewayDomainID returns a new GatewayDomainId struct
func NewGatewayDomainID(subscriptionId string, resourceGroupName string, springName string, gatewayName string, domainName string) GatewayDomainId {
	return GatewayDomainId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		GatewayName:       gatewayName,
		DomainName:        domainName,
	}
}

// ParseGatewayDomainID parses 'input' into a GatewayDomainId
func ParseGatewayDomainID(input string) (*GatewayDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(GatewayDomainId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GatewayDomainId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.GatewayName, ok = parsed.Parsed["gatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "gatewayName", *parsed)
	}

	if id.DomainName, ok = parsed.Parsed["domainName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "domainName", *parsed)
	}

	return &id, nil
}

// ParseGatewayDomainIDInsensitively parses 'input' case-insensitively into a GatewayDomainId
// note: this method should only be used for API response data and not user input
func ParseGatewayDomainIDInsensitively(input string) (*GatewayDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(GatewayDomainId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GatewayDomainId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.GatewayName, ok = parsed.Parsed["gatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "gatewayName", *parsed)
	}

	if id.DomainName, ok = parsed.Parsed["domainName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "domainName", *parsed)
	}

	return &id, nil
}

// ValidateGatewayDomainID checks that 'input' can be parsed as a Gateway Domain ID
func ValidateGatewayDomainID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGatewayDomainID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Gateway Domain ID
func (id GatewayDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/gateways/%s/domains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.GatewayName, id.DomainName)
}

// Segments returns a slice of Resource ID Segments which comprise this Gateway Domain ID
func (id GatewayDomainId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springValue"),
		resourceids.StaticSegment("staticGateways", "gateways", "gateways"),
		resourceids.UserSpecifiedSegment("gatewayName", "gatewayValue"),
		resourceids.StaticSegment("staticDomains", "domains", "domains"),
		resourceids.UserSpecifiedSegment("domainName", "domainValue"),
	}
}

// String returns a human-readable description of this Gateway Domain ID
func (id GatewayDomainId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Gateway Name: %q", id.GatewayName),
		fmt.Sprintf("Domain Name: %q", id.DomainName),
	}
	return fmt.Sprintf("Gateway Domain (%s)", strings.Join(components, "\n"))
}
