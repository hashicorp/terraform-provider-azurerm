package testlines

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CommunicationsGatewayId{})
}

var _ resourceids.ResourceId = &CommunicationsGatewayId{}

// CommunicationsGatewayId is a struct representing the Resource ID for a Communications Gateway
type CommunicationsGatewayId struct {
	SubscriptionId            string
	ResourceGroupName         string
	CommunicationsGatewayName string
}

// NewCommunicationsGatewayID returns a new CommunicationsGatewayId struct
func NewCommunicationsGatewayID(subscriptionId string, resourceGroupName string, communicationsGatewayName string) CommunicationsGatewayId {
	return CommunicationsGatewayId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		CommunicationsGatewayName: communicationsGatewayName,
	}
}

// ParseCommunicationsGatewayID parses 'input' into a CommunicationsGatewayId
func ParseCommunicationsGatewayID(input string) (*CommunicationsGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CommunicationsGatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CommunicationsGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCommunicationsGatewayIDInsensitively parses 'input' case-insensitively into a CommunicationsGatewayId
// note: this method should only be used for API response data and not user input
func ParseCommunicationsGatewayIDInsensitively(input string) (*CommunicationsGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CommunicationsGatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CommunicationsGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CommunicationsGatewayId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CommunicationsGatewayName, ok = input.Parsed["communicationsGatewayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "communicationsGatewayName", input)
	}

	return nil
}

// ValidateCommunicationsGatewayID checks that 'input' can be parsed as a Communications Gateway ID
func ValidateCommunicationsGatewayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCommunicationsGatewayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Communications Gateway ID
func (id CommunicationsGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.VoiceServices/communicationsGateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CommunicationsGatewayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Communications Gateway ID
func (id CommunicationsGatewayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftVoiceServices", "Microsoft.VoiceServices", "Microsoft.VoiceServices"),
		resourceids.StaticSegment("staticCommunicationsGateways", "communicationsGateways", "communicationsGateways"),
		resourceids.UserSpecifiedSegment("communicationsGatewayName", "communicationsGatewayName"),
	}
}

// String returns a human-readable description of this Communications Gateway ID
func (id CommunicationsGatewayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Communications Gateway Name: %q", id.CommunicationsGatewayName),
	}
	return fmt.Sprintf("Communications Gateway (%s)", strings.Join(components, "\n"))
}
