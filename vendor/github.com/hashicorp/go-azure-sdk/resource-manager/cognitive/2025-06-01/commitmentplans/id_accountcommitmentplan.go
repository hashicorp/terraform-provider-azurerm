package commitmentplans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AccountCommitmentPlanId{})
}

var _ resourceids.ResourceId = &AccountCommitmentPlanId{}

// AccountCommitmentPlanId is a struct representing the Resource ID for a Account Commitment Plan
type AccountCommitmentPlanId struct {
	SubscriptionId     string
	ResourceGroupName  string
	AccountName        string
	CommitmentPlanName string
}

// NewAccountCommitmentPlanID returns a new AccountCommitmentPlanId struct
func NewAccountCommitmentPlanID(subscriptionId string, resourceGroupName string, accountName string, commitmentPlanName string) AccountCommitmentPlanId {
	return AccountCommitmentPlanId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		AccountName:        accountName,
		CommitmentPlanName: commitmentPlanName,
	}
}

// ParseAccountCommitmentPlanID parses 'input' into a AccountCommitmentPlanId
func ParseAccountCommitmentPlanID(input string) (*AccountCommitmentPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccountCommitmentPlanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccountCommitmentPlanId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAccountCommitmentPlanIDInsensitively parses 'input' case-insensitively into a AccountCommitmentPlanId
// note: this method should only be used for API response data and not user input
func ParseAccountCommitmentPlanIDInsensitively(input string) (*AccountCommitmentPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccountCommitmentPlanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccountCommitmentPlanId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AccountCommitmentPlanId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.CommitmentPlanName, ok = input.Parsed["commitmentPlanName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "commitmentPlanName", input)
	}

	return nil
}

// ValidateAccountCommitmentPlanID checks that 'input' can be parsed as a Account Commitment Plan ID
func ValidateAccountCommitmentPlanID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAccountCommitmentPlanID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Account Commitment Plan ID
func (id AccountCommitmentPlanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/accounts/%s/commitmentPlans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.CommitmentPlanName)
}

// Segments returns a slice of Resource ID Segments which comprise this Account Commitment Plan ID
func (id AccountCommitmentPlanId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCognitiveServices", "Microsoft.CognitiveServices", "Microsoft.CognitiveServices"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountName"),
		resourceids.StaticSegment("staticCommitmentPlans", "commitmentPlans", "commitmentPlans"),
		resourceids.UserSpecifiedSegment("commitmentPlanName", "commitmentPlanName"),
	}
}

// String returns a human-readable description of this Account Commitment Plan ID
func (id AccountCommitmentPlanId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Commitment Plan Name: %q", id.CommitmentPlanName),
	}
	return fmt.Sprintf("Account Commitment Plan (%s)", strings.Join(components, "\n"))
}
