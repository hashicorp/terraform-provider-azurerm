package prefixlistlocalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LocalRuleStackPrefixListId{}

// LocalRuleStackPrefixListId is a struct representing the Resource ID for a Local Rule Stack Prefix List
type LocalRuleStackPrefixListId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LocalRuleStackName string
	PrefixListName     string
}

// NewLocalRuleStackPrefixListID returns a new LocalRuleStackPrefixListId struct
func NewLocalRuleStackPrefixListID(subscriptionId string, resourceGroupName string, localRuleStackName string, prefixListName string) LocalRuleStackPrefixListId {
	return LocalRuleStackPrefixListId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LocalRuleStackName: localRuleStackName,
		PrefixListName:     prefixListName,
	}
}

// ParseLocalRuleStackPrefixListID parses 'input' into a LocalRuleStackPrefixListId
func ParseLocalRuleStackPrefixListID(input string) (*LocalRuleStackPrefixListId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRuleStackPrefixListId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRuleStackPrefixListId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRuleStackName, ok = parsed.Parsed["localRuleStackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRuleStackName", *parsed)
	}

	if id.PrefixListName, ok = parsed.Parsed["prefixListName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "prefixListName", *parsed)
	}

	return &id, nil
}

// ParseLocalRuleStackPrefixListIDInsensitively parses 'input' case-insensitively into a LocalRuleStackPrefixListId
// note: this method should only be used for API response data and not user input
func ParseLocalRuleStackPrefixListIDInsensitively(input string) (*LocalRuleStackPrefixListId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRuleStackPrefixListId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRuleStackPrefixListId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRuleStackName, ok = parsed.Parsed["localRuleStackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRuleStackName", *parsed)
	}

	if id.PrefixListName, ok = parsed.Parsed["prefixListName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "prefixListName", *parsed)
	}

	return &id, nil
}

// ValidateLocalRuleStackPrefixListID checks that 'input' can be parsed as a Local Rule Stack Prefix List ID
func ValidateLocalRuleStackPrefixListID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocalRuleStackPrefixListID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Local Rule Stack Prefix List ID
func (id LocalRuleStackPrefixListId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/PaloAltoNetworks.CloudNGFW/localRuleStacks/%s/prefixLists/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocalRuleStackName, id.PrefixListName)
}

// Segments returns a slice of Resource ID Segments which comprise this Local Rule Stack Prefix List ID
func (id LocalRuleStackPrefixListId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudNGFW", "PaloAltoNetworks.CloudNGFW", "PaloAltoNetworks.CloudNGFW"),
		resourceids.StaticSegment("staticLocalRuleStacks", "localRuleStacks", "localRuleStacks"),
		resourceids.UserSpecifiedSegment("localRuleStackName", "localRuleStackValue"),
		resourceids.StaticSegment("staticPrefixLists", "prefixLists", "prefixLists"),
		resourceids.UserSpecifiedSegment("prefixListName", "prefixListValue"),
	}
}

// String returns a human-readable description of this Local Rule Stack Prefix List ID
func (id LocalRuleStackPrefixListId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Local Rule Stack Name: %q", id.LocalRuleStackName),
		fmt.Sprintf("Prefix List Name: %q", id.PrefixListName),
	}
	return fmt.Sprintf("Local Rule Stack Prefix List (%s)", strings.Join(components, "\n"))
}
