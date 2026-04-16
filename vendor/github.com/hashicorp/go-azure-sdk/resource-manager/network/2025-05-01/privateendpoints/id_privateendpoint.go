package privateendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PrivateEndpointId{})
}

var _ resourceids.ResourceId = &PrivateEndpointId{}

// PrivateEndpointId is a struct representing the Resource ID for a Private Endpoint
type PrivateEndpointId struct {
	SubscriptionId      string
	ResourceGroupName   string
	PrivateEndpointName string
}

// NewPrivateEndpointID returns a new PrivateEndpointId struct
func NewPrivateEndpointID(subscriptionId string, resourceGroupName string, privateEndpointName string) PrivateEndpointId {
	return PrivateEndpointId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		PrivateEndpointName: privateEndpointName,
	}
}

// ParsePrivateEndpointID parses 'input' into a PrivateEndpointId
func ParsePrivateEndpointID(input string) (*PrivateEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePrivateEndpointIDInsensitively parses 'input' case-insensitively into a PrivateEndpointId
// note: this method should only be used for API response data and not user input
func ParsePrivateEndpointIDInsensitively(input string) (*PrivateEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PrivateEndpointId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PrivateEndpointName, ok = input.Parsed["privateEndpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointName", input)
	}

	return nil
}

// ValidatePrivateEndpointID checks that 'input' can be parsed as a Private Endpoint ID
func ValidatePrivateEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Endpoint ID
func (id PrivateEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Endpoint ID
func (id PrivateEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticPrivateEndpoints", "privateEndpoints", "privateEndpoints"),
		resourceids.UserSpecifiedSegment("privateEndpointName", "privateEndpointName"),
	}
}

// String returns a human-readable description of this Private Endpoint ID
func (id PrivateEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Endpoint Name: %q", id.PrivateEndpointName),
	}
	return fmt.Sprintf("Private Endpoint (%s)", strings.Join(components, "\n"))
}
