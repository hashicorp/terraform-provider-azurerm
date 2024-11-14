package raipolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RaiPolicyId{})
}

var _ resourceids.ResourceId = &RaiPolicyId{}

// RaiPolicyId is a struct representing the Resource ID for a Rai Policy
type RaiPolicyId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	RaiPolicyName     string
}

// NewRaiPolicyID returns a new RaiPolicyId struct
func NewRaiPolicyID(subscriptionId string, resourceGroupName string, accountName string, raiPolicyName string) RaiPolicyId {
	return RaiPolicyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		RaiPolicyName:     raiPolicyName,
	}
}

// ParseRaiPolicyID parses 'input' into a RaiPolicyId
func ParseRaiPolicyID(input string) (*RaiPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RaiPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RaiPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRaiPolicyIDInsensitively parses 'input' case-insensitively into a RaiPolicyId
// note: this method should only be used for API response data and not user input
func ParseRaiPolicyIDInsensitively(input string) (*RaiPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RaiPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RaiPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RaiPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AccountName, ok = input.Parsed["accountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "accountName", input)
	}

	if id.RaiPolicyName, ok = input.Parsed["raiPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "raiPolicyName", input)
	}

	return nil
}

// ValidateRaiPolicyID checks that 'input' can be parsed as a Rai Policy ID
func ValidateRaiPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRaiPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Rai Policy ID
func (id RaiPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/accounts/%s/raiPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.RaiPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Rai Policy ID
func (id RaiPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCognitiveServices", "Microsoft.CognitiveServices", "Microsoft.CognitiveServices"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountName"),
		resourceids.StaticSegment("staticRaiPolicies", "raiPolicies", "raiPolicies"),
		resourceids.UserSpecifiedSegment("raiPolicyName", "raiPolicyName"),
	}
}

// String returns a human-readable description of this Rai Policy ID
func (id RaiPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Rai Policy Name: %q", id.RaiPolicyName),
	}
	return fmt.Sprintf("Rai Policy (%s)", strings.Join(components, "\n"))
}
