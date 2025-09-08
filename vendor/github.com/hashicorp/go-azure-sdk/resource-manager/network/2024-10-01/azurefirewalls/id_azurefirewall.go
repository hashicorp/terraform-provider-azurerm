package azurefirewalls

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AzureFirewallId{})
}

var _ resourceids.ResourceId = &AzureFirewallId{}

// AzureFirewallId is a struct representing the Resource ID for a Azure Firewall
type AzureFirewallId struct {
	SubscriptionId    string
	ResourceGroupName string
	AzureFirewallName string
}

// NewAzureFirewallID returns a new AzureFirewallId struct
func NewAzureFirewallID(subscriptionId string, resourceGroupName string, azureFirewallName string) AzureFirewallId {
	return AzureFirewallId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AzureFirewallName: azureFirewallName,
	}
}

// ParseAzureFirewallID parses 'input' into a AzureFirewallId
func ParseAzureFirewallID(input string) (*AzureFirewallId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AzureFirewallId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AzureFirewallId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAzureFirewallIDInsensitively parses 'input' case-insensitively into a AzureFirewallId
// note: this method should only be used for API response data and not user input
func ParseAzureFirewallIDInsensitively(input string) (*AzureFirewallId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AzureFirewallId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AzureFirewallId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AzureFirewallId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AzureFirewallName, ok = input.Parsed["azureFirewallName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "azureFirewallName", input)
	}

	return nil
}

// ValidateAzureFirewallID checks that 'input' can be parsed as a Azure Firewall ID
func ValidateAzureFirewallID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAzureFirewallID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Azure Firewall ID
func (id AzureFirewallId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/azureFirewalls/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AzureFirewallName)
}

// Segments returns a slice of Resource ID Segments which comprise this Azure Firewall ID
func (id AzureFirewallId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticAzureFirewalls", "azureFirewalls", "azureFirewalls"),
		resourceids.UserSpecifiedSegment("azureFirewallName", "azureFirewallName"),
	}
}

// String returns a human-readable description of this Azure Firewall ID
func (id AzureFirewallId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Azure Firewall Name: %q", id.AzureFirewallName),
	}
	return fmt.Sprintf("Azure Firewall (%s)", strings.Join(components, "\n"))
}
