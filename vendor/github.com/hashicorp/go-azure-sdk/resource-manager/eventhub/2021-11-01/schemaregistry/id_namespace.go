package schemaregistry

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NamespaceId{})
}

var _ resourceids.ResourceId = &NamespaceId{}

// NamespaceId is a struct representing the Resource ID for a Namespace
type NamespaceId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
}

// NewNamespaceID returns a new NamespaceId struct
func NewNamespaceID(subscriptionId string, resourceGroupName string, namespaceName string) NamespaceId {
	return NamespaceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
	}
}

// ParseNamespaceID parses 'input' into a NamespaceId
func ParseNamespaceID(input string) (*NamespaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NamespaceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NamespaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNamespaceIDInsensitively parses 'input' case-insensitively into a NamespaceId
// note: this method should only be used for API response data and not user input
func ParseNamespaceIDInsensitively(input string) (*NamespaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NamespaceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NamespaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NamespaceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NamespaceName, ok = input.Parsed["namespaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", input)
	}

	return nil
}

// ValidateNamespaceID checks that 'input' can be parsed as a Namespace ID
func ValidateNamespaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNamespaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Namespace ID
func (id NamespaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Namespace ID
func (id NamespaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventHub", "Microsoft.EventHub", "Microsoft.EventHub"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
	}
}

// String returns a human-readable description of this Namespace ID
func (id NamespaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
	}
	return fmt.Sprintf("Namespace (%s)", strings.Join(components, "\n"))
}
