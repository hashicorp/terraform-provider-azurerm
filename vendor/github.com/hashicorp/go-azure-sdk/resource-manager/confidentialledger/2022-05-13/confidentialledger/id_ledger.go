package confidentialledger

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LedgerId{}

// LedgerId is a struct representing the Resource ID for a Ledger
type LedgerId struct {
	SubscriptionId    string
	ResourceGroupName string
	LedgerName        string
}

// NewLedgerID returns a new LedgerId struct
func NewLedgerID(subscriptionId string, resourceGroupName string, ledgerName string) LedgerId {
	return LedgerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LedgerName:        ledgerName,
	}
}

// ParseLedgerID parses 'input' into a LedgerId
func ParseLedgerID(input string) (*LedgerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LedgerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LedgerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LedgerName, ok = parsed.Parsed["ledgerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ledgerName", *parsed)
	}

	return &id, nil
}

// ParseLedgerIDInsensitively parses 'input' case-insensitively into a LedgerId
// note: this method should only be used for API response data and not user input
func ParseLedgerIDInsensitively(input string) (*LedgerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LedgerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LedgerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LedgerName, ok = parsed.Parsed["ledgerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ledgerName", *parsed)
	}

	return &id, nil
}

// ValidateLedgerID checks that 'input' can be parsed as a Ledger ID
func ValidateLedgerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLedgerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Ledger ID
func (id LedgerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ConfidentialLedger/ledgers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LedgerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Ledger ID
func (id LedgerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftConfidentialLedger", "Microsoft.ConfidentialLedger", "Microsoft.ConfidentialLedger"),
		resourceids.StaticSegment("staticLedgers", "ledgers", "ledgers"),
		resourceids.UserSpecifiedSegment("ledgerName", "ledgerValue"),
	}
}

// String returns a human-readable description of this Ledger ID
func (id LedgerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Ledger Name: %q", id.LedgerName),
	}
	return fmt.Sprintf("Ledger (%s)", strings.Join(components, "\n"))
}
