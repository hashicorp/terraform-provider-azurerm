package prefixlistlocalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LocalRulestackPrefixListId{})
}

var _ resourceids.ResourceId = &LocalRulestackPrefixListId{}

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
	parser := resourceids.NewParserFromResourceIdType(&LocalRulestackPrefixListId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocalRulestackPrefixListId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLocalRulestackPrefixListIDInsensitively parses 'input' case-insensitively into a LocalRulestackPrefixListId
// note: this method should only be used for API response data and not user input
func ParseLocalRulestackPrefixListIDInsensitively(input string) (*LocalRulestackPrefixListId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocalRulestackPrefixListId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocalRulestackPrefixListId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LocalRulestackPrefixListId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.LocalRulestackName, ok = input.Parsed["localRulestackName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "localRulestackName", input)
	}

	if id.PrefixListName, ok = input.Parsed["prefixListName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "prefixListName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("localRulestackName", "localRulestackName"),
		resourceids.StaticSegment("staticPrefixLists", "prefixLists", "prefixLists"),
		resourceids.UserSpecifiedSegment("prefixListName", "prefixListName"),
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
