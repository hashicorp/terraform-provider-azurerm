package routes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AfdEndpointId{})
}

var _ resourceids.ResourceId = &AfdEndpointId{}

// AfdEndpointId is a struct representing the Resource ID for a Afd Endpoint
type AfdEndpointId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	AfdEndpointName   string
}

// NewAfdEndpointID returns a new AfdEndpointId struct
func NewAfdEndpointID(subscriptionId string, resourceGroupName string, profileName string, afdEndpointName string) AfdEndpointId {
	return AfdEndpointId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		AfdEndpointName:   afdEndpointName,
	}
}

// ParseAfdEndpointID parses 'input' into a AfdEndpointId
func ParseAfdEndpointID(input string) (*AfdEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AfdEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AfdEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAfdEndpointIDInsensitively parses 'input' case-insensitively into a AfdEndpointId
// note: this method should only be used for API response data and not user input
func ParseAfdEndpointIDInsensitively(input string) (*AfdEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AfdEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AfdEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AfdEndpointId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AfdEndpointName, ok = input.Parsed["afdEndpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "afdEndpointName", input)
	}

	return nil
}

// ValidateAfdEndpointID checks that 'input' can be parsed as a Afd Endpoint ID
func ValidateAfdEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAfdEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Afd Endpoint ID
func (id AfdEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/afdEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.AfdEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Afd Endpoint ID
func (id AfdEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCdn", "Microsoft.Cdn", "Microsoft.Cdn"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileName"),
		resourceids.StaticSegment("staticAfdEndpoints", "afdEndpoints", "afdEndpoints"),
		resourceids.UserSpecifiedSegment("afdEndpointName", "afdEndpointName"),
	}
}

// String returns a human-readable description of this Afd Endpoint ID
func (id AfdEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Afd Endpoint Name: %q", id.AfdEndpointName),
	}
	return fmt.Sprintf("Afd Endpoint (%s)", strings.Join(components, "\n"))
}
