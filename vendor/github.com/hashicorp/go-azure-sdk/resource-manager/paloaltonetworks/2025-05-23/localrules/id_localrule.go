package localrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LocalRuleId{})
}

var _ resourceids.ResourceId = &LocalRuleId{}

// LocalRuleId is a struct representing the Resource ID for a Local Rule
type LocalRuleId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LocalRulestackName string
	LocalRuleName      string
}

// NewLocalRuleID returns a new LocalRuleId struct
func NewLocalRuleID(subscriptionId string, resourceGroupName string, localRulestackName string, localRuleName string) LocalRuleId {
	return LocalRuleId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LocalRulestackName: localRulestackName,
		LocalRuleName:      localRuleName,
	}
}

// ParseLocalRuleID parses 'input' into a LocalRuleId
func ParseLocalRuleID(input string) (*LocalRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocalRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocalRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLocalRuleIDInsensitively parses 'input' case-insensitively into a LocalRuleId
// note: this method should only be used for API response data and not user input
func ParseLocalRuleIDInsensitively(input string) (*LocalRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LocalRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LocalRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LocalRuleId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.LocalRuleName, ok = input.Parsed["localRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "localRuleName", input)
	}

	return nil
}

// ValidateLocalRuleID checks that 'input' can be parsed as a Local Rule ID
func ValidateLocalRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocalRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Local Rule ID
func (id LocalRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/%s/localRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName, id.LocalRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Local Rule ID
func (id LocalRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudngfw", "PaloAltoNetworks.Cloudngfw", "PaloAltoNetworks.Cloudngfw"),
		resourceids.StaticSegment("staticLocalRulestacks", "localRulestacks", "localRulestacks"),
		resourceids.UserSpecifiedSegment("localRulestackName", "localRulestackName"),
		resourceids.StaticSegment("staticLocalRules", "localRules", "localRules"),
		resourceids.UserSpecifiedSegment("localRuleName", "localRuleName"),
	}
}

// String returns a human-readable description of this Local Rule ID
func (id LocalRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Local Rulestack Name: %q", id.LocalRulestackName),
		fmt.Sprintf("Local Rule Name: %q", id.LocalRuleName),
	}
	return fmt.Sprintf("Local Rule (%s)", strings.Join(components, "\n"))
}
