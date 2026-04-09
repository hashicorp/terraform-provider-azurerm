package integrationaccountagreements

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&IntegrationAccountId{})
}

var _ resourceids.ResourceId = &IntegrationAccountId{}

// IntegrationAccountId is a struct representing the Resource ID for a Integration Account
type IntegrationAccountId struct {
	SubscriptionId         string
	ResourceGroupName      string
	IntegrationAccountName string
}

// NewIntegrationAccountID returns a new IntegrationAccountId struct
func NewIntegrationAccountID(subscriptionId string, resourceGroupName string, integrationAccountName string) IntegrationAccountId {
	return IntegrationAccountId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		IntegrationAccountName: integrationAccountName,
	}
}

// ParseIntegrationAccountID parses 'input' into a IntegrationAccountId
func ParseIntegrationAccountID(input string) (*IntegrationAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IntegrationAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IntegrationAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseIntegrationAccountIDInsensitively parses 'input' case-insensitively into a IntegrationAccountId
// note: this method should only be used for API response data and not user input
func ParseIntegrationAccountIDInsensitively(input string) (*IntegrationAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IntegrationAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IntegrationAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *IntegrationAccountId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.IntegrationAccountName, ok = input.Parsed["integrationAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "integrationAccountName", input)
	}

	return nil
}

// ValidateIntegrationAccountID checks that 'input' can be parsed as a Integration Account ID
func ValidateIntegrationAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIntegrationAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Integration Account ID
func (id IntegrationAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IntegrationAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Integration Account ID
func (id IntegrationAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticIntegrationAccounts", "integrationAccounts", "integrationAccounts"),
		resourceids.UserSpecifiedSegment("integrationAccountName", "integrationAccountName"),
	}
}

// String returns a human-readable description of this Integration Account ID
func (id IntegrationAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Integration Account Name: %q", id.IntegrationAccountName),
	}
	return fmt.Sprintf("Integration Account (%s)", strings.Join(components, "\n"))
}
