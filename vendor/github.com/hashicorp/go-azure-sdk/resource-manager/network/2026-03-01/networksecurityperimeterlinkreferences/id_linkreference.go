package networksecurityperimeterlinkreferences

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LinkReferenceId{})
}

var _ resourceids.ResourceId = &LinkReferenceId{}

// LinkReferenceId is a struct representing the Resource ID for a Link Reference
type LinkReferenceId struct {
	SubscriptionId               string
	ResourceGroupName            string
	NetworkSecurityPerimeterName string
	LinkReferenceName            string
}

// NewLinkReferenceID returns a new LinkReferenceId struct
func NewLinkReferenceID(subscriptionId string, resourceGroupName string, networkSecurityPerimeterName string, linkReferenceName string) LinkReferenceId {
	return LinkReferenceId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		NetworkSecurityPerimeterName: networkSecurityPerimeterName,
		LinkReferenceName:            linkReferenceName,
	}
}

// ParseLinkReferenceID parses 'input' into a LinkReferenceId
func ParseLinkReferenceID(input string) (*LinkReferenceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LinkReferenceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LinkReferenceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLinkReferenceIDInsensitively parses 'input' case-insensitively into a LinkReferenceId
// note: this method should only be used for API response data and not user input
func ParseLinkReferenceIDInsensitively(input string) (*LinkReferenceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LinkReferenceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LinkReferenceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LinkReferenceId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.LinkReferenceName, ok = input.Parsed["linkReferenceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "linkReferenceName", input)
	}

	return nil
}

// ValidateLinkReferenceID checks that 'input' can be parsed as a Link Reference ID
func ValidateLinkReferenceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLinkReferenceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Link Reference ID
func (id LinkReferenceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityPerimeters/%s/linkReferences/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityPerimeterName, id.LinkReferenceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Link Reference ID
func (id LinkReferenceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkSecurityPerimeters", "networkSecurityPerimeters", "networkSecurityPerimeters"),
		resourceids.UserSpecifiedSegment("networkSecurityPerimeterName", "networkSecurityPerimeterName"),
		resourceids.StaticSegment("staticLinkReferences", "linkReferences", "linkReferences"),
		resourceids.UserSpecifiedSegment("linkReferenceName", "linkReferenceName"),
	}
}

// String returns a human-readable description of this Link Reference ID
func (id LinkReferenceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Security Perimeter Name: %q", id.NetworkSecurityPerimeterName),
		fmt.Sprintf("Link Reference Name: %q", id.LinkReferenceName),
	}
	return fmt.Sprintf("Link Reference (%s)", strings.Join(components, "\n"))
}
