package networksecurityperimeterassociations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ResourceAssociationId{})
}

var _ resourceids.ResourceId = &ResourceAssociationId{}

// ResourceAssociationId is a struct representing the Resource ID for a Resource Association
type ResourceAssociationId struct {
	SubscriptionId               string
	ResourceGroupName            string
	NetworkSecurityPerimeterName string
	ResourceAssociationName      string
}

// NewResourceAssociationID returns a new ResourceAssociationId struct
func NewResourceAssociationID(subscriptionId string, resourceGroupName string, networkSecurityPerimeterName string, resourceAssociationName string) ResourceAssociationId {
	return ResourceAssociationId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		NetworkSecurityPerimeterName: networkSecurityPerimeterName,
		ResourceAssociationName:      resourceAssociationName,
	}
}

// ParseResourceAssociationID parses 'input' into a ResourceAssociationId
func ParseResourceAssociationID(input string) (*ResourceAssociationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceAssociationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceAssociationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseResourceAssociationIDInsensitively parses 'input' case-insensitively into a ResourceAssociationId
// note: this method should only be used for API response data and not user input
func ParseResourceAssociationIDInsensitively(input string) (*ResourceAssociationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceAssociationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceAssociationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ResourceAssociationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkSecurityPerimeterName, ok = input.Parsed["networkSecurityPerimeterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkSecurityPerimeterName", input)
	}

	if id.ResourceAssociationName, ok = input.Parsed["resourceAssociationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceAssociationName", input)
	}

	return nil
}

// ValidateResourceAssociationID checks that 'input' can be parsed as a Resource Association ID
func ValidateResourceAssociationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResourceAssociationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resource Association ID
func (id ResourceAssociationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityPerimeters/%s/resourceAssociations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityPerimeterName, id.ResourceAssociationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Resource Association ID
func (id ResourceAssociationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkSecurityPerimeters", "networkSecurityPerimeters", "networkSecurityPerimeters"),
		resourceids.UserSpecifiedSegment("networkSecurityPerimeterName", "networkSecurityPerimeterName"),
		resourceids.StaticSegment("staticResourceAssociations", "resourceAssociations", "resourceAssociations"),
		resourceids.UserSpecifiedSegment("resourceAssociationName", "resourceAssociationName"),
	}
}

// String returns a human-readable description of this Resource Association ID
func (id ResourceAssociationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Security Perimeter Name: %q", id.NetworkSecurityPerimeterName),
		fmt.Sprintf("Resource Association Name: %q", id.ResourceAssociationName),
	}
	return fmt.Sprintf("Resource Association (%s)", strings.Join(components, "\n"))
}
