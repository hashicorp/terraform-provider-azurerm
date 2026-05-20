package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GatewayDomainId{})
}

var _ resourceids.ResourceId = &GatewayDomainId{}

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
	parser := resourceids.NewParserFromResourceIdType(&GatewayDomainId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GatewayDomainId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGatewayDomainIDInsensitively parses 'input' case-insensitively into a GatewayDomainId
// note: this method should only be used for API response data and not user input
func ParseGatewayDomainIDInsensitively(input string) (*GatewayDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GatewayDomainId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GatewayDomainId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GatewayDomainId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpringName, ok = input.Parsed["springName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "springName", input)
	}

	if id.GatewayName, ok = input.Parsed["gatewayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "gatewayName", input)
	}

	if id.DomainName, ok = input.Parsed["domainName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "domainName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticGateways", "gateways", "gateways"),
		resourceids.UserSpecifiedSegment("gatewayName", "gatewayName"),
		resourceids.StaticSegment("staticDomains", "domains", "domains"),
		resourceids.UserSpecifiedSegment("domainName", "domainName"),
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
