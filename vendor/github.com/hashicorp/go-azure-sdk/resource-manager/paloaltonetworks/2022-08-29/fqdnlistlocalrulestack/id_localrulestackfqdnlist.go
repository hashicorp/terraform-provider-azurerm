package fqdnlistlocalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LocalRulestackFqdnListId{}

// LocalRulestackFqdnListId is a struct representing the Resource ID for a Local Rulestack Fqdn List
type LocalRulestackFqdnListId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LocalRulestackName string
	FqdnListName       string
}

// NewLocalRulestackFqdnListID returns a new LocalRulestackFqdnListId struct
func NewLocalRulestackFqdnListID(subscriptionId string, resourceGroupName string, localRulestackName string, fqdnListName string) LocalRulestackFqdnListId {
	return LocalRulestackFqdnListId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LocalRulestackName: localRulestackName,
		FqdnListName:       fqdnListName,
	}
}

// ParseLocalRulestackFqdnListID parses 'input' into a LocalRulestackFqdnListId
func ParseLocalRulestackFqdnListID(input string) (*LocalRulestackFqdnListId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRulestackFqdnListId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRulestackFqdnListId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRulestackName, ok = parsed.Parsed["localRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRulestackName", *parsed)
	}

	if id.FqdnListName, ok = parsed.Parsed["fqdnListName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fqdnListName", *parsed)
	}

	return &id, nil
}

// ParseLocalRulestackFqdnListIDInsensitively parses 'input' case-insensitively into a LocalRulestackFqdnListId
// note: this method should only be used for API response data and not user input
func ParseLocalRulestackFqdnListIDInsensitively(input string) (*LocalRulestackFqdnListId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRulestackFqdnListId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRulestackFqdnListId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRulestackName, ok = parsed.Parsed["localRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRulestackName", *parsed)
	}

	if id.FqdnListName, ok = parsed.Parsed["fqdnListName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fqdnListName", *parsed)
	}

	return &id, nil
}

// ValidateLocalRulestackFqdnListID checks that 'input' can be parsed as a Local Rulestack Fqdn List ID
func ValidateLocalRulestackFqdnListID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocalRulestackFqdnListID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Local Rulestack Fqdn List ID
func (id LocalRulestackFqdnListId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/%s/fqdnLists/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName, id.FqdnListName)
}

// Segments returns a slice of Resource ID Segments which comprise this Local Rulestack Fqdn List ID
func (id LocalRulestackFqdnListId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudngfw", "PaloAltoNetworks.Cloudngfw", "PaloAltoNetworks.Cloudngfw"),
		resourceids.StaticSegment("staticLocalRulestacks", "localRulestacks", "localRulestacks"),
		resourceids.UserSpecifiedSegment("localRulestackName", "localRulestackValue"),
		resourceids.StaticSegment("staticFqdnLists", "fqdnLists", "fqdnLists"),
		resourceids.UserSpecifiedSegment("fqdnListName", "fqdnListValue"),
	}
}

// String returns a human-readable description of this Local Rulestack Fqdn List ID
func (id LocalRulestackFqdnListId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Local Rulestack Name: %q", id.LocalRulestackName),
		fmt.Sprintf("Fqdn List Name: %q", id.FqdnListName),
	}
	return fmt.Sprintf("Local Rulestack Fqdn List (%s)", strings.Join(components, "\n"))
}
