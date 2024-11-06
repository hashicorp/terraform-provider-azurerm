package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&HybridConnectionId{})
}

var _ resourceids.ResourceId = &HybridConnectionId{}

// HybridConnectionId is a struct representing the Resource ID for a Hybrid Connection
type HybridConnectionId struct {
	SubscriptionId       string
	ResourceGroupName    string
	SiteName             string
	HybridConnectionName string
}

// NewHybridConnectionID returns a new HybridConnectionId struct
func NewHybridConnectionID(subscriptionId string, resourceGroupName string, siteName string, hybridConnectionName string) HybridConnectionId {
	return HybridConnectionId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		SiteName:             siteName,
		HybridConnectionName: hybridConnectionName,
	}
}

// ParseHybridConnectionID parses 'input' into a HybridConnectionId
func ParseHybridConnectionID(input string) (*HybridConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HybridConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HybridConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHybridConnectionIDInsensitively parses 'input' case-insensitively into a HybridConnectionId
// note: this method should only be used for API response data and not user input
func ParseHybridConnectionIDInsensitively(input string) (*HybridConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HybridConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HybridConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HybridConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.HybridConnectionName, ok = input.Parsed["hybridConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hybridConnectionName", input)
	}

	return nil
}

// ValidateHybridConnectionID checks that 'input' can be parsed as a Hybrid Connection ID
func ValidateHybridConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHybridConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hybrid Connection ID
func (id HybridConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/hybridConnection/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.HybridConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hybrid Connection ID
func (id HybridConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticHybridConnection", "hybridConnection", "hybridConnection"),
		resourceids.UserSpecifiedSegment("hybridConnectionName", "hybridConnectionName"),
	}
}

// String returns a human-readable description of this Hybrid Connection ID
func (id HybridConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Hybrid Connection Name: %q", id.HybridConnectionName),
	}
	return fmt.Sprintf("Hybrid Connection (%s)", strings.Join(components, "\n"))
}
