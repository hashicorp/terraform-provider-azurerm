package prefixlistlocalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LocalRulestackPrefixListId{}

// LocalRulestackPrefixListId is a struct representing the Resource ID for a Local Rulestack Prefix List
type LocalRulestackPrefixListId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LocalRulestackName string
	PrefixListName     string
}

// NewLocalRulestackPrefixListID returns a new LocalRulestackPrefixListId struct
func NewLocalRulestackPrefixListID(subscriptionId string, resourceGroupName string, localRulestackName string, prefixListName string) LocalRulestackPrefixListId {
	return LocalRulestackPrefixListId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LocalRulestackName: localRulestackName,
		PrefixListName:     prefixListName,
	}
}

// ParseLocalRulestackPrefixListID parses 'input' into a LocalRulestackPrefixListId
func ParseLocalRulestackPrefixListID(input string) (*LocalRulestackPrefixListId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRulestackPrefixListId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRulestackPrefixListId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRulestackName, ok = parsed.Parsed["localRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRulestackName", *parsed)
	}

	if id.PrefixListName, ok = parsed.Parsed["prefixListName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "prefixListName", *parsed)
	}

	return &id, nil
}

// ParseLocalRulestackPrefixListIDInsensitively parses 'input' case-insensitively into a LocalRulestackPrefixListId
// note: this method should only be used for API response data and not user input
func ParseLocalRulestackPrefixListIDInsensitively(input string) (*LocalRulestackPrefixListId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRulestackPrefixListId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRulestackPrefixListId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRulestackName, ok = parsed.Parsed["localRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRulestackName", *parsed)
	}

	if id.PrefixListName, ok = parsed.Parsed["prefixListName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "prefixListName", *parsed)
	}

	return &id, nil
}

// ValidateLocalRulestackPrefixListID checks that 'input' can be parsed as a Local Rulestack Prefix List ID
func ValidateLocalRulestackPrefixListID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocalRulestackPrefixListID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Local Rulestack Prefix List ID
func (id LocalRulestackPrefixListId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/%s/prefixLists/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName, id.PrefixListName)
}

// Segments returns a slice of Resource ID Segments which comprise this Local Rulestack Prefix List ID
func (id LocalRulestackPrefixListId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudngfw", "PaloAltoNetworks.Cloudngfw", "PaloAltoNetworks.Cloudngfw"),
		resourceids.StaticSegment("staticLocalRulestacks", "localRulestacks", "localRulestacks"),
		resourceids.UserSpecifiedSegment("localRulestackName", "localRulestackValue"),
		resourceids.StaticSegment("staticPrefixLists", "prefixLists", "prefixLists"),
		resourceids.UserSpecifiedSegment("prefixListName", "prefixListValue"),
	}
}

// String returns a human-readable description of this Local Rulestack Prefix List ID
func (id LocalRulestackPrefixListId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Local Rulestack Name: %q", id.LocalRulestackName),
		fmt.Sprintf("Prefix List Name: %q", id.PrefixListName),
	}
	return fmt.Sprintf("Local Rulestack Prefix List (%s)", strings.Join(components, "\n"))
}
