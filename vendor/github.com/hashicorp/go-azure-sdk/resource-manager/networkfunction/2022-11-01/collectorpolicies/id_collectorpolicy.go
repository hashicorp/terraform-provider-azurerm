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
	recaser.RegisterResourceId(&CollectorPolicyId{})
}

var _ resourceids.ResourceId = &CollectorPolicyId{}

// CollectorPolicyId is a struct representing the Resource ID for a Collector Policy
type CollectorPolicyId struct {
	SubscriptionId            string
	ResourceGroupName         string
	AzureTrafficCollectorName string
	CollectorPolicyName       string
}

// NewCollectorPolicyID returns a new CollectorPolicyId struct
func NewCollectorPolicyID(subscriptionId string, resourceGroupName string, azureTrafficCollectorName string, collectorPolicyName string) CollectorPolicyId {
	return CollectorPolicyId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		AzureTrafficCollectorName: azureTrafficCollectorName,
		CollectorPolicyName:       collectorPolicyName,
	}
}

// ParseCollectorPolicyID parses 'input' into a CollectorPolicyId
func ParseCollectorPolicyID(input string) (*CollectorPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CollectorPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CollectorPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCollectorPolicyIDInsensitively parses 'input' case-insensitively into a CollectorPolicyId
// note: this method should only be used for API response data and not user input
func ParseCollectorPolicyIDInsensitively(input string) (*CollectorPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CollectorPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CollectorPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CollectorPolicyId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.CollectorPolicyName, ok = input.Parsed["collectorPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "collectorPolicyName", input)
	}

	return nil
}

// ValidateCollectorPolicyID checks that 'input' can be parsed as a Collector Policy ID
func ValidateCollectorPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCollectorPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Collector Policy ID
func (id CollectorPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetworkFunction/azureTrafficCollectors/%s/collectorPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AzureTrafficCollectorName, id.CollectorPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Collector Policy ID
func (id CollectorPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetworkFunction", "Microsoft.NetworkFunction", "Microsoft.NetworkFunction"),
		resourceids.StaticSegment("staticAzureTrafficCollectors", "azureTrafficCollectors", "azureTrafficCollectors"),
		resourceids.UserSpecifiedSegment("azureTrafficCollectorName", "azureTrafficCollectorName"),
		resourceids.StaticSegment("staticCollectorPolicies", "collectorPolicies", "collectorPolicies"),
		resourceids.UserSpecifiedSegment("collectorPolicyName", "collectorPolicyName"),
	}
}

// String returns a human-readable description of this Collector Policy ID
func (id CollectorPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Azure Traffic Collector Name: %q", id.AzureTrafficCollectorName),
		fmt.Sprintf("Collector Policy Name: %q", id.CollectorPolicyName),
	}
	return fmt.Sprintf("Collector Policy (%s)", strings.Join(components, "\n"))
}
