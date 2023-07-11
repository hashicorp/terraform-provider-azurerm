package fqdnlistlocalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LocalRuleStackFqdnListId{}

// LocalRuleStackFqdnListId is a struct representing the Resource ID for a Local Rule Stack Fqdn List
type LocalRuleStackFqdnListId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LocalRuleStackName string
	FqdnListName       string
}

// NewLocalRuleStackFqdnListID returns a new LocalRuleStackFqdnListId struct
func NewLocalRuleStackFqdnListID(subscriptionId string, resourceGroupName string, localRuleStackName string, fqdnListName string) LocalRuleStackFqdnListId {
	return LocalRuleStackFqdnListId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LocalRuleStackName: localRuleStackName,
		FqdnListName:       fqdnListName,
	}
}

// ParseLocalRuleStackFqdnListID parses 'input' into a LocalRuleStackFqdnListId
func ParseLocalRuleStackFqdnListID(input string) (*LocalRuleStackFqdnListId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRuleStackFqdnListId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRuleStackFqdnListId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRuleStackName, ok = parsed.Parsed["localRuleStackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRuleStackName", *parsed)
	}

	if id.FqdnListName, ok = parsed.Parsed["fqdnListName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fqdnListName", *parsed)
	}

	return &id, nil
}

// ParseLocalRuleStackFqdnListIDInsensitively parses 'input' case-insensitively into a LocalRuleStackFqdnListId
// note: this method should only be used for API response data and not user input
func ParseLocalRuleStackFqdnListIDInsensitively(input string) (*LocalRuleStackFqdnListId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRuleStackFqdnListId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRuleStackFqdnListId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRuleStackName, ok = parsed.Parsed["localRuleStackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRuleStackName", *parsed)
	}

	if id.FqdnListName, ok = parsed.Parsed["fqdnListName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fqdnListName", *parsed)
	}

	return &id, nil
}

// ValidateLocalRuleStackFqdnListID checks that 'input' can be parsed as a Local Rule Stack Fqdn List ID
func ValidateLocalRuleStackFqdnListID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocalRuleStackFqdnListID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Local Rule Stack Fqdn List ID
func (id LocalRuleStackFqdnListId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/PaloAltoNetworks.CloudNGFW/localRuleStacks/%s/fqdnLists/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocalRuleStackName, id.FqdnListName)
}

// Segments returns a slice of Resource ID Segments which comprise this Local Rule Stack Fqdn List ID
func (id LocalRuleStackFqdnListId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudNGFW", "PaloAltoNetworks.CloudNGFW", "PaloAltoNetworks.CloudNGFW"),
		resourceids.StaticSegment("staticLocalRuleStacks", "localRuleStacks", "localRuleStacks"),
		resourceids.UserSpecifiedSegment("localRuleStackName", "localRuleStackValue"),
		resourceids.StaticSegment("staticFqdnLists", "fqdnLists", "fqdnLists"),
		resourceids.UserSpecifiedSegment("fqdnListName", "fqdnListValue"),
	}
}

// String returns a human-readable description of this Local Rule Stack Fqdn List ID
func (id LocalRuleStackFqdnListId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Local Rule Stack Name: %q", id.LocalRuleStackName),
		fmt.Sprintf("Fqdn List Name: %q", id.FqdnListName),
	}
	return fmt.Sprintf("Local Rule Stack Fqdn List (%s)", strings.Join(components, "\n"))
}
