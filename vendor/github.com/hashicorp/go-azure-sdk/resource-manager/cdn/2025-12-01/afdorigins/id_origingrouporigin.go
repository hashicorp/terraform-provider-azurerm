package afdorigins

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OriginGroupOriginId{})
}

var _ resourceids.ResourceId = &OriginGroupOriginId{}

// OriginGroupOriginId is a struct representing the Resource ID for a Origin Group Origin
type OriginGroupOriginId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	OriginGroupName   string
	OriginName        string
}

// NewOriginGroupOriginID returns a new OriginGroupOriginId struct
func NewOriginGroupOriginID(subscriptionId string, resourceGroupName string, profileName string, originGroupName string, originName string) OriginGroupOriginId {
	return OriginGroupOriginId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		OriginGroupName:   originGroupName,
		OriginName:        originName,
	}
}

// ParseOriginGroupOriginID parses 'input' into a OriginGroupOriginId
func ParseOriginGroupOriginID(input string) (*OriginGroupOriginId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OriginGroupOriginId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OriginGroupOriginId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOriginGroupOriginIDInsensitively parses 'input' case-insensitively into a OriginGroupOriginId
// note: this method should only be used for API response data and not user input
func ParseOriginGroupOriginIDInsensitively(input string) (*OriginGroupOriginId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OriginGroupOriginId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OriginGroupOriginId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OriginGroupOriginId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.OriginGroupName, ok = input.Parsed["originGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "originGroupName", input)
	}

	if id.OriginName, ok = input.Parsed["originName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "originName", input)
	}

	return nil
}

// ValidateOriginGroupOriginID checks that 'input' can be parsed as a Origin Group Origin ID
func ValidateOriginGroupOriginID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOriginGroupOriginID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Origin Group Origin ID
func (id OriginGroupOriginId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/originGroups/%s/origins/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.OriginGroupName, id.OriginName)
}

// Segments returns a slice of Resource ID Segments which comprise this Origin Group Origin ID
func (id OriginGroupOriginId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCdn", "Microsoft.Cdn", "Microsoft.Cdn"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileName"),
		resourceids.StaticSegment("staticOriginGroups", "originGroups", "originGroups"),
		resourceids.UserSpecifiedSegment("originGroupName", "originGroupName"),
		resourceids.StaticSegment("staticOrigins", "origins", "origins"),
		resourceids.UserSpecifiedSegment("originName", "originName"),
	}
}

// String returns a human-readable description of this Origin Group Origin ID
func (id OriginGroupOriginId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Origin Group Name: %q", id.OriginGroupName),
		fmt.Sprintf("Origin Name: %q", id.OriginName),
	}
	return fmt.Sprintf("Origin Group Origin (%s)", strings.Join(components, "\n"))
}
