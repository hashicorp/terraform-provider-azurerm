package prefixlistlocalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LocalRuleStackId{}

// LocalRuleStackId is a struct representing the Resource ID for a Local Rule Stack
type LocalRuleStackId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LocalRuleStackName string
}

// NewLocalRuleStackID returns a new LocalRuleStackId struct
func NewLocalRuleStackID(subscriptionId string, resourceGroupName string, localRuleStackName string) LocalRuleStackId {
	return LocalRuleStackId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LocalRuleStackName: localRuleStackName,
	}
}

// ParseLocalRuleStackID parses 'input' into a LocalRuleStackId
func ParseLocalRuleStackID(input string) (*LocalRuleStackId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRuleStackId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRuleStackId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRuleStackName, ok = parsed.Parsed["localRuleStackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRuleStackName", *parsed)
	}

	return &id, nil
}

// ParseLocalRuleStackIDInsensitively parses 'input' case-insensitively into a LocalRuleStackId
// note: this method should only be used for API response data and not user input
func ParseLocalRuleStackIDInsensitively(input string) (*LocalRuleStackId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRuleStackId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRuleStackId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRuleStackName, ok = parsed.Parsed["localRuleStackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRuleStackName", *parsed)
	}

	return &id, nil
}

// ValidateLocalRuleStackID checks that 'input' can be parsed as a Local Rule Stack ID
func ValidateLocalRuleStackID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocalRuleStackID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Local Rule Stack ID
func (id LocalRuleStackId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/PaloAltoNetworks.CloudNGFW/localRuleStacks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocalRuleStackName)
}

// Segments returns a slice of Resource ID Segments which comprise this Local Rule Stack ID
func (id LocalRuleStackId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudNGFW", "PaloAltoNetworks.CloudNGFW", "PaloAltoNetworks.CloudNGFW"),
		resourceids.StaticSegment("staticLocalRuleStacks", "localRuleStacks", "localRuleStacks"),
		resourceids.UserSpecifiedSegment("localRuleStackName", "localRuleStackValue"),
	}
}

// String returns a human-readable description of this Local Rule Stack ID
func (id LocalRuleStackId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Local Rule Stack Name: %q", id.LocalRuleStackName),
	}
	return fmt.Sprintf("Local Rule Stack (%s)", strings.Join(components, "\n"))
}
