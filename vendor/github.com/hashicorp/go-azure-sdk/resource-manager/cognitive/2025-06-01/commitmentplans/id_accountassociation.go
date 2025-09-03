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
	recaser.RegisterResourceId(&AccountAssociationId{})
}

var _ resourceids.ResourceId = &AccountAssociationId{}

// AccountAssociationId is a struct representing the Resource ID for a Account Association
type AccountAssociationId struct {
	SubscriptionId         string
	ResourceGroupName      string
	CommitmentPlanName     string
	AccountAssociationName string
}

// NewAccountAssociationID returns a new AccountAssociationId struct
func NewAccountAssociationID(subscriptionId string, resourceGroupName string, commitmentPlanName string, accountAssociationName string) AccountAssociationId {
	return AccountAssociationId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		CommitmentPlanName:     commitmentPlanName,
		AccountAssociationName: accountAssociationName,
	}
}

// ParseAccountAssociationID parses 'input' into a AccountAssociationId
func ParseAccountAssociationID(input string) (*AccountAssociationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccountAssociationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccountAssociationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAccountAssociationIDInsensitively parses 'input' case-insensitively into a AccountAssociationId
// note: this method should only be used for API response data and not user input
func ParseAccountAssociationIDInsensitively(input string) (*AccountAssociationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccountAssociationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccountAssociationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AccountAssociationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CommitmentPlanName, ok = input.Parsed["commitmentPlanName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "commitmentPlanName", input)
	}

	if id.AccountAssociationName, ok = input.Parsed["accountAssociationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "accountAssociationName", input)
	}

	return nil
}

// ValidateAccountAssociationID checks that 'input' can be parsed as a Account Association ID
func ValidateAccountAssociationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAccountAssociationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Account Association ID
func (id AccountAssociationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/commitmentPlans/%s/accountAssociations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CommitmentPlanName, id.AccountAssociationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Account Association ID
func (id AccountAssociationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCognitiveServices", "Microsoft.CognitiveServices", "Microsoft.CognitiveServices"),
		resourceids.StaticSegment("staticCommitmentPlans", "commitmentPlans", "commitmentPlans"),
		resourceids.UserSpecifiedSegment("commitmentPlanName", "commitmentPlanName"),
		resourceids.StaticSegment("staticAccountAssociations", "accountAssociations", "accountAssociations"),
		resourceids.UserSpecifiedSegment("accountAssociationName", "accountAssociationName"),
	}
}

// String returns a human-readable description of this Account Association ID
func (id AccountAssociationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Commitment Plan Name: %q", id.CommitmentPlanName),
		fmt.Sprintf("Account Association Name: %q", id.AccountAssociationName),
	}
	return fmt.Sprintf("Account Association (%s)", strings.Join(components, "\n"))
}
