package localrulestacks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LocalRulestackId{})
}

var _ resourceids.ResourceId = &LocalRulestackId{}

// LocalRulestackId is a struct representing the Resource ID for a Local Rulestack
type LocalRulestackId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LocalRulestackName string
}

// NewLocalRulestackID returns a new LocalRulestackId struct
func NewLocalRulestackID(subscriptionId string, resourceGroupName string, localRulestackName string) LocalRulestackId {
	return LocalRulestackId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LocalRulestackName: localRulestackName,
	}
}

// ParseLocalRulestackID parses 'input' into a LocalRulestackId
func ParseLocalRulestackID(input string) (*LocalRulestackId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocalRulestackId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocalRulestackId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLocalRulestackIDInsensitively parses 'input' case-insensitively into a LocalRulestackId
// note: this method should only be used for API response data and not user input
func ParseLocalRulestackIDInsensitively(input string) (*LocalRulestackId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocalRulestackId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocalRulestackId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LocalRulestackId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateLocalRulestackID checks that 'input' can be parsed as a Local Rulestack ID
func ValidateLocalRulestackID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocalRulestackID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Local Rulestack ID
func (id LocalRulestackId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
}

// Segments returns a slice of Resource ID Segments which comprise this Local Rulestack ID
func (id LocalRulestackId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudngfw", "PaloAltoNetworks.Cloudngfw", "PaloAltoNetworks.Cloudngfw"),
		resourceids.StaticSegment("staticLocalRulestacks", "localRulestacks", "localRulestacks"),
		resourceids.UserSpecifiedSegment("localRulestackName", "localRulestackName"),
	}
}

// String returns a human-readable description of this Local Rulestack ID
func (id LocalRulestackId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Local Rulestack Name: %q", id.LocalRulestackName),
	}
	return fmt.Sprintf("Local Rulestack (%s)", strings.Join(components, "\n"))
}
