package clients

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ClientId{})
}

var _ resourceids.ResourceId = &ClientId{}

// ClientId is a struct representing the Resource ID for a Client
type ClientId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	ClientName        string
}

// NewClientID returns a new ClientId struct
func NewClientID(subscriptionId string, resourceGroupName string, namespaceName string, clientName string) ClientId {
	return ClientId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		ClientName:        clientName,
	}
}

// ParseClientID parses 'input' into a ClientId
func ParseClientID(input string) (*ClientId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ClientId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ClientId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseClientIDInsensitively parses 'input' case-insensitively into a ClientId
// note: this method should only be used for API response data and not user input
func ParseClientIDInsensitively(input string) (*ClientId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ClientId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ClientId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ClientId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ClientName, ok = input.Parsed["clientName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "clientName", input)
	}

	return nil
}

// ValidateClientID checks that 'input' can be parsed as a Client ID
func ValidateClientID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseClientID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Client ID
func (id ClientId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/namespaces/%s/clients/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.ClientName)
}

// Segments returns a slice of Resource ID Segments which comprise this Client ID
func (id ClientId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticClients", "clients", "clients"),
		resourceids.UserSpecifiedSegment("clientName", "clientName"),
	}
}

// String returns a human-readable description of this Client ID
func (id ClientId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Client Name: %q", id.ClientName),
	}
	return fmt.Sprintf("Client (%s)", strings.Join(components, "\n"))
}
