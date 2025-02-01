package origins

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OriginId{})
}

var _ resourceids.ResourceId = &OriginId{}

// OriginId is a struct representing the Resource ID for a Origin
type OriginId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	EndpointName      string
	OriginName        string
}

// NewOriginID returns a new OriginId struct
func NewOriginID(subscriptionId string, resourceGroupName string, profileName string, endpointName string, originName string) OriginId {
	return OriginId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		EndpointName:      endpointName,
		OriginName:        originName,
	}
}

// ParseOriginID parses 'input' into a OriginId
func ParseOriginID(input string) (*OriginId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OriginId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OriginId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOriginIDInsensitively parses 'input' case-insensitively into a OriginId
// note: this method should only be used for API response data and not user input
func ParseOriginIDInsensitively(input string) (*OriginId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OriginId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OriginId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OriginId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ProfileName, ok = input.Parsed["profileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "profileName", input)
	}

	if id.EndpointName, ok = input.Parsed["endpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "endpointName", input)
	}

	if id.OriginName, ok = input.Parsed["originName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "originName", input)
	}

	return nil
}

// ValidateOriginID checks that 'input' can be parsed as a Origin ID
func ValidateOriginID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOriginID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Origin ID
func (id OriginId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/endpoints/%s/origins/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.OriginName)
}

// Segments returns a slice of Resource ID Segments which comprise this Origin ID
func (id OriginId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCdn", "Microsoft.Cdn", "Microsoft.Cdn"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileName"),
		resourceids.StaticSegment("staticEndpoints", "endpoints", "endpoints"),
		resourceids.UserSpecifiedSegment("endpointName", "endpointName"),
		resourceids.StaticSegment("staticOrigins", "origins", "origins"),
		resourceids.UserSpecifiedSegment("originName", "originName"),
	}
}

// String returns a human-readable description of this Origin ID
func (id OriginId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Endpoint Name: %q", id.EndpointName),
		fmt.Sprintf("Origin Name: %q", id.OriginName),
	}
	return fmt.Sprintf("Origin (%s)", strings.Join(components, "\n"))
}
