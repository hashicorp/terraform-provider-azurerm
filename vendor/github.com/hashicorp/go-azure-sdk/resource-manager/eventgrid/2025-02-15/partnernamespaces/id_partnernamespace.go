package partnernamespaces

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PartnerNamespaceId{})
}

var _ resourceids.ResourceId = &PartnerNamespaceId{}

// PartnerNamespaceId is a struct representing the Resource ID for a Partner Namespace
type PartnerNamespaceId struct {
	SubscriptionId       string
	ResourceGroupName    string
	PartnerNamespaceName string
}

// NewPartnerNamespaceID returns a new PartnerNamespaceId struct
func NewPartnerNamespaceID(subscriptionId string, resourceGroupName string, partnerNamespaceName string) PartnerNamespaceId {
	return PartnerNamespaceId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		PartnerNamespaceName: partnerNamespaceName,
	}
}

// ParsePartnerNamespaceID parses 'input' into a PartnerNamespaceId
func ParsePartnerNamespaceID(input string) (*PartnerNamespaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PartnerNamespaceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PartnerNamespaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePartnerNamespaceIDInsensitively parses 'input' case-insensitively into a PartnerNamespaceId
// note: this method should only be used for API response data and not user input
func ParsePartnerNamespaceIDInsensitively(input string) (*PartnerNamespaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PartnerNamespaceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PartnerNamespaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PartnerNamespaceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PartnerNamespaceName, ok = input.Parsed["partnerNamespaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "partnerNamespaceName", input)
	}

	return nil
}

// ValidatePartnerNamespaceID checks that 'input' can be parsed as a Partner Namespace ID
func ValidatePartnerNamespaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePartnerNamespaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Partner Namespace ID
func (id PartnerNamespaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/partnerNamespaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PartnerNamespaceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Partner Namespace ID
func (id PartnerNamespaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticPartnerNamespaces", "partnerNamespaces", "partnerNamespaces"),
		resourceids.UserSpecifiedSegment("partnerNamespaceName", "partnerNamespaceName"),
	}
}

// String returns a human-readable description of this Partner Namespace ID
func (id PartnerNamespaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Partner Namespace Name: %q", id.PartnerNamespaceName),
	}
	return fmt.Sprintf("Partner Namespace (%s)", strings.Join(components, "\n"))
}
