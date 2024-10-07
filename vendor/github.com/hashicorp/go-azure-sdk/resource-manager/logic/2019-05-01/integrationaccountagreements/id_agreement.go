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
	recaser.RegisterResourceId(&AgreementId{})
}

var _ resourceids.ResourceId = &AgreementId{}

// AgreementId is a struct representing the Resource ID for a Agreement
type AgreementId struct {
	SubscriptionId         string
	ResourceGroupName      string
	IntegrationAccountName string
	AgreementName          string
}

// NewAgreementID returns a new AgreementId struct
func NewAgreementID(subscriptionId string, resourceGroupName string, integrationAccountName string, agreementName string) AgreementId {
	return AgreementId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		IntegrationAccountName: integrationAccountName,
		AgreementName:          agreementName,
	}
}

// ParseAgreementID parses 'input' into a AgreementId
func ParseAgreementID(input string) (*AgreementId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AgreementId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AgreementId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAgreementIDInsensitively parses 'input' case-insensitively into a AgreementId
// note: this method should only be used for API response data and not user input
func ParseAgreementIDInsensitively(input string) (*AgreementId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AgreementId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AgreementId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AgreementId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AgreementName, ok = input.Parsed["agreementName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "agreementName", input)
	}

	return nil
}

// ValidateAgreementID checks that 'input' can be parsed as a Agreement ID
func ValidateAgreementID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAgreementID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Agreement ID
func (id AgreementId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/agreements/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IntegrationAccountName, id.AgreementName)
}

// Segments returns a slice of Resource ID Segments which comprise this Agreement ID
func (id AgreementId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticIntegrationAccounts", "integrationAccounts", "integrationAccounts"),
		resourceids.UserSpecifiedSegment("integrationAccountName", "integrationAccountName"),
		resourceids.StaticSegment("staticAgreements", "agreements", "agreements"),
		resourceids.UserSpecifiedSegment("agreementName", "agreementName"),
	}
}

// String returns a human-readable description of this Agreement ID
func (id AgreementId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Integration Account Name: %q", id.IntegrationAccountName),
		fmt.Sprintf("Agreement Name: %q", id.AgreementName),
	}
	return fmt.Sprintf("Agreement (%s)", strings.Join(components, "\n"))
}
