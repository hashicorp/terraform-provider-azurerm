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
	recaser.RegisterResourceId(&AccountRaiPolicyId{})
}

var _ resourceids.ResourceId = &AccountRaiPolicyId{}

// AccountRaiPolicyId is a struct representing the Resource ID for a Account Rai Policy
type AccountRaiPolicyId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	RaiPolicyName     string
}

// NewAccountRaiPolicyID returns a new AccountRaiPolicyId struct
func NewAccountRaiPolicyID(subscriptionId string, resourceGroupName string, accountName string, raiPolicyName string) AccountRaiPolicyId {
	return AccountRaiPolicyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		RaiPolicyName:     raiPolicyName,
	}
}

// ParseAccountRaiPolicyID parses 'input' into a AccountRaiPolicyId
func ParseAccountRaiPolicyID(input string) (*AccountRaiPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccountRaiPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccountRaiPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAccountRaiPolicyIDInsensitively parses 'input' case-insensitively into a AccountRaiPolicyId
// note: this method should only be used for API response data and not user input
func ParseAccountRaiPolicyIDInsensitively(input string) (*AccountRaiPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccountRaiPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccountRaiPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AccountRaiPolicyId) FromParseResult(input resourceids.ParseResult) error {
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

// ValidateAccountRaiPolicyID checks that 'input' can be parsed as a Account Rai Policy ID
func ValidateAccountRaiPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAccountRaiPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Account Rai Policy ID
func (id AccountRaiPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/accounts/%s/raiPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.RaiPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Account Rai Policy ID
func (id AccountRaiPolicyId) Segments() []resourceids.Segment {
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

// String returns a human-readable description of this Account Rai Policy ID
func (id AccountRaiPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Rai Policy Name: %q", id.RaiPolicyName),
	}
	return fmt.Sprintf("Account Rai Policy (%s)", strings.Join(components, "\n"))
}
