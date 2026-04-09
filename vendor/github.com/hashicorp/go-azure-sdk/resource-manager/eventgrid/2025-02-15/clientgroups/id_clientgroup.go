package clientgroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ClientGroupId{})
}

var _ resourceids.ResourceId = &ClientGroupId{}

// ClientGroupId is a struct representing the Resource ID for a Client Group
type ClientGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	ClientGroupName   string
}

// NewClientGroupID returns a new ClientGroupId struct
func NewClientGroupID(subscriptionId string, resourceGroupName string, namespaceName string, clientGroupName string) ClientGroupId {
	return ClientGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		ClientGroupName:   clientGroupName,
	}
}

// ParseClientGroupID parses 'input' into a ClientGroupId
func ParseClientGroupID(input string) (*ClientGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ClientGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ClientGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseClientGroupIDInsensitively parses 'input' case-insensitively into a ClientGroupId
// note: this method should only be used for API response data and not user input
func ParseClientGroupIDInsensitively(input string) (*ClientGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ClientGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ClientGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ClientGroupId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ClientGroupName, ok = input.Parsed["clientGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "clientGroupName", input)
	}

	return nil
}

// ValidateClientGroupID checks that 'input' can be parsed as a Client Group ID
func ValidateClientGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseClientGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Client Group ID
func (id ClientGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/namespaces/%s/clientGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.ClientGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Client Group ID
func (id ClientGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticClientGroups", "clientGroups", "clientGroups"),
		resourceids.UserSpecifiedSegment("clientGroupName", "clientGroupName"),
	}
}

// String returns a human-readable description of this Client Group ID
func (id ClientGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Client Group Name: %q", id.ClientGroupName),
	}
	return fmt.Sprintf("Client Group (%s)", strings.Join(components, "\n"))
}
