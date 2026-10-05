package collectorpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AzureTrafficCollectorId{})
}

var _ resourceids.ResourceId = &AzureTrafficCollectorId{}

// AzureTrafficCollectorId is a struct representing the Resource ID for a Azure Traffic Collector
type AzureTrafficCollectorId struct {
	SubscriptionId            string
	ResourceGroupName         string
	AzureTrafficCollectorName string
}

// NewAzureTrafficCollectorID returns a new AzureTrafficCollectorId struct
func NewAzureTrafficCollectorID(subscriptionId string, resourceGroupName string, azureTrafficCollectorName string) AzureTrafficCollectorId {
	return AzureTrafficCollectorId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		AzureTrafficCollectorName: azureTrafficCollectorName,
	}
}

// ParseAzureTrafficCollectorID parses 'input' into a AzureTrafficCollectorId
func ParseAzureTrafficCollectorID(input string) (*AzureTrafficCollectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AzureTrafficCollectorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AzureTrafficCollectorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAzureTrafficCollectorIDInsensitively parses 'input' case-insensitively into a AzureTrafficCollectorId
// note: this method should only be used for API response data and not user input
func ParseAzureTrafficCollectorIDInsensitively(input string) (*AzureTrafficCollectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AzureTrafficCollectorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AzureTrafficCollectorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AzureTrafficCollectorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AzureTrafficCollectorName, ok = input.Parsed["azureTrafficCollectorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "azureTrafficCollectorName", input)
	}

	return nil
}

// ValidateAzureTrafficCollectorID checks that 'input' can be parsed as a Azure Traffic Collector ID
func ValidateAzureTrafficCollectorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAzureTrafficCollectorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Azure Traffic Collector ID
func (id AzureTrafficCollectorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetworkFunction/azureTrafficCollectors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AzureTrafficCollectorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Azure Traffic Collector ID
func (id AzureTrafficCollectorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetworkFunction", "Microsoft.NetworkFunction", "Microsoft.NetworkFunction"),
		resourceids.StaticSegment("staticAzureTrafficCollectors", "azureTrafficCollectors", "azureTrafficCollectors"),
		resourceids.UserSpecifiedSegment("azureTrafficCollectorName", "azureTrafficCollectorName"),
	}
}

// String returns a human-readable description of this Azure Traffic Collector ID
func (id AzureTrafficCollectorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Azure Traffic Collector Name: %q", id.AzureTrafficCollectorName),
	}
	return fmt.Sprintf("Azure Traffic Collector (%s)", strings.Join(components, "\n"))
}
