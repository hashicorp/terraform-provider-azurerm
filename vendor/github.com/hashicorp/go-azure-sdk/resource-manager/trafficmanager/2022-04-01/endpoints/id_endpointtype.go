package endpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&EndpointTypeId{})
}

var _ resourceids.ResourceId = &EndpointTypeId{}

// EndpointTypeId is a struct representing the Resource ID for a Endpoint Type
type EndpointTypeId struct {
	SubscriptionId            string
	ResourceGroupName         string
	TrafficManagerProfileName string
	EndpointType              EndpointType
	EndpointName              string
}

// NewEndpointTypeID returns a new EndpointTypeId struct
func NewEndpointTypeID(subscriptionId string, resourceGroupName string, trafficManagerProfileName string, endpointType EndpointType, endpointName string) EndpointTypeId {
	return EndpointTypeId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		TrafficManagerProfileName: trafficManagerProfileName,
		EndpointType:              endpointType,
		EndpointName:              endpointName,
	}
}

// ParseEndpointTypeID parses 'input' into a EndpointTypeId
func ParseEndpointTypeID(input string) (*EndpointTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EndpointTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EndpointTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseEndpointTypeIDInsensitively parses 'input' case-insensitively into a EndpointTypeId
// note: this method should only be used for API response data and not user input
func ParseEndpointTypeIDInsensitively(input string) (*EndpointTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EndpointTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EndpointTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *EndpointTypeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.TrafficManagerProfileName, ok = input.Parsed["trafficManagerProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "trafficManagerProfileName", input)
	}

	if v, ok := input.Parsed["endpointType"]; true {
		if !ok {
			return resourceids.NewSegmentNotSpecifiedError(id, "endpointType", input)
		}

		endpointType, err := parseEndpointType(v)
		if err != nil {
			return fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.EndpointType = *endpointType
	}

	if id.EndpointName, ok = input.Parsed["endpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "endpointName", input)
	}

	return nil
}

// ValidateEndpointTypeID checks that 'input' can be parsed as a Endpoint Type ID
func ValidateEndpointTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEndpointTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Endpoint Type ID
func (id EndpointTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/trafficManagerProfiles/%s/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.TrafficManagerProfileName, string(id.EndpointType), id.EndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Endpoint Type ID
func (id EndpointTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticTrafficManagerProfiles", "trafficManagerProfiles", "trafficManagerProfiles"),
		resourceids.UserSpecifiedSegment("trafficManagerProfileName", "trafficManagerProfileName"),
		resourceids.ConstantSegment("endpointType", PossibleValuesForEndpointType(), "AzureEndpoints"),
		resourceids.UserSpecifiedSegment("endpointName", "endpointName"),
	}
}

// String returns a human-readable description of this Endpoint Type ID
func (id EndpointTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Traffic Manager Profile Name: %q", id.TrafficManagerProfileName),
		fmt.Sprintf("Endpoint Type: %q", string(id.EndpointType)),
		fmt.Sprintf("Endpoint Name: %q", id.EndpointName),
	}
	return fmt.Sprintf("Endpoint Type (%s)", strings.Join(components, "\n"))
}
