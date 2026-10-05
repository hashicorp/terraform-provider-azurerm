package staticsites

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StaticSitePrivateEndpointConnectionId{})
}

var _ resourceids.ResourceId = &StaticSitePrivateEndpointConnectionId{}

// StaticSitePrivateEndpointConnectionId is a struct representing the Resource ID for a Static Site Private Endpoint Connection
type StaticSitePrivateEndpointConnectionId struct {
	SubscriptionId                string
	ResourceGroupName             string
	StaticSiteName                string
	PrivateEndpointConnectionName string
}

// NewStaticSitePrivateEndpointConnectionID returns a new StaticSitePrivateEndpointConnectionId struct
func NewStaticSitePrivateEndpointConnectionID(subscriptionId string, resourceGroupName string, staticSiteName string, privateEndpointConnectionName string) StaticSitePrivateEndpointConnectionId {
	return StaticSitePrivateEndpointConnectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		StaticSiteName:                staticSiteName,
		PrivateEndpointConnectionName: privateEndpointConnectionName,
	}
}

// ParseStaticSitePrivateEndpointConnectionID parses 'input' into a StaticSitePrivateEndpointConnectionId
func ParseStaticSitePrivateEndpointConnectionID(input string) (*StaticSitePrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StaticSitePrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StaticSitePrivateEndpointConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStaticSitePrivateEndpointConnectionIDInsensitively parses 'input' case-insensitively into a StaticSitePrivateEndpointConnectionId
// note: this method should only be used for API response data and not user input
func ParseStaticSitePrivateEndpointConnectionIDInsensitively(input string) (*StaticSitePrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StaticSitePrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StaticSitePrivateEndpointConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StaticSitePrivateEndpointConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StaticSiteName, ok = input.Parsed["staticSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "staticSiteName", input)
	}

	if id.PrivateEndpointConnectionName, ok = input.Parsed["privateEndpointConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointConnectionName", input)
	}

	return nil
}

// ValidateStaticSitePrivateEndpointConnectionID checks that 'input' can be parsed as a Static Site Private Endpoint Connection ID
func ValidateStaticSitePrivateEndpointConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStaticSitePrivateEndpointConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Static Site Private Endpoint Connection ID
func (id StaticSitePrivateEndpointConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/staticSites/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName, id.PrivateEndpointConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Static Site Private Endpoint Connection ID
func (id StaticSitePrivateEndpointConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticStaticSites", "staticSites", "staticSites"),
		resourceids.UserSpecifiedSegment("staticSiteName", "staticSiteName"),
		resourceids.StaticSegment("staticPrivateEndpointConnections", "privateEndpointConnections", "privateEndpointConnections"),
		resourceids.UserSpecifiedSegment("privateEndpointConnectionName", "privateEndpointConnectionName"),
	}
}

// String returns a human-readable description of this Static Site Private Endpoint Connection ID
func (id StaticSitePrivateEndpointConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Static Site Name: %q", id.StaticSiteName),
		fmt.Sprintf("Private Endpoint Connection Name: %q", id.PrivateEndpointConnectionName),
	}
	return fmt.Sprintf("Static Site Private Endpoint Connection (%s)", strings.Join(components, "\n"))
}
