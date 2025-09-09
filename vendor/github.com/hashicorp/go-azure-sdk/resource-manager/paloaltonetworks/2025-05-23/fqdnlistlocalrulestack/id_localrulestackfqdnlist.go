package fqdnlistlocalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LocalRulestackFqdnListId{})
}

var _ resourceids.ResourceId = &LocalRulestackFqdnListId{}

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
	parser := resourceids.NewParserFromResourceIdType(&LocalRulestackFqdnListId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocalRulestackFqdnListId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLocalRulestackFqdnListIDInsensitively parses 'input' case-insensitively into a LocalRulestackFqdnListId
// note: this method should only be used for API response data and not user input
func ParseLocalRulestackFqdnListIDInsensitively(input string) (*LocalRulestackFqdnListId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocalRulestackFqdnListId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocalRulestackFqdnListId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LocalRulestackFqdnListId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.FqdnListName, ok = input.Parsed["fqdnListName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "fqdnListName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("localRulestackName", "localRulestackName"),
		resourceids.StaticSegment("staticFqdnLists", "fqdnLists", "fqdnLists"),
		resourceids.UserSpecifiedSegment("fqdnListName", "fqdnListName"),
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
