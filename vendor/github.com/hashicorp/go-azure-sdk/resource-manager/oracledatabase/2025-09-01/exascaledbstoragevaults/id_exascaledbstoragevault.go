package exascaledbstoragevaults

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExascaleDbStorageVaultId{})
}

var _ resourceids.ResourceId = &ExascaleDbStorageVaultId{}

// ExascaleDbStorageVaultId is a struct representing the Resource ID for a Exascale Db Storage Vault
type ExascaleDbStorageVaultId struct {
	SubscriptionId             string
	ResourceGroupName          string
	ExascaleDbStorageVaultName string
}

// NewExascaleDbStorageVaultID returns a new ExascaleDbStorageVaultId struct
func NewExascaleDbStorageVaultID(subscriptionId string, resourceGroupName string, exascaleDbStorageVaultName string) ExascaleDbStorageVaultId {
	return ExascaleDbStorageVaultId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		ExascaleDbStorageVaultName: exascaleDbStorageVaultName,
	}
}

// ParseExascaleDbStorageVaultID parses 'input' into a ExascaleDbStorageVaultId
func ParseExascaleDbStorageVaultID(input string) (*ExascaleDbStorageVaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExascaleDbStorageVaultId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExascaleDbStorageVaultId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExascaleDbStorageVaultIDInsensitively parses 'input' case-insensitively into a ExascaleDbStorageVaultId
// note: this method should only be used for API response data and not user input
func ParseExascaleDbStorageVaultIDInsensitively(input string) (*ExascaleDbStorageVaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExascaleDbStorageVaultId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExascaleDbStorageVaultId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExascaleDbStorageVaultId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExascaleDbStorageVaultName, ok = input.Parsed["exascaleDbStorageVaultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "exascaleDbStorageVaultName", input)
	}

	return nil
}

// ValidateExascaleDbStorageVaultID checks that 'input' can be parsed as a Exascale Db Storage Vault ID
func ValidateExascaleDbStorageVaultID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExascaleDbStorageVaultID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Exascale Db Storage Vault ID
func (id ExascaleDbStorageVaultId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/exascaleDbStorageVaults/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExascaleDbStorageVaultName)
}

// Segments returns a slice of Resource ID Segments which comprise this Exascale Db Storage Vault ID
func (id ExascaleDbStorageVaultId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticExascaleDbStorageVaults", "exascaleDbStorageVaults", "exascaleDbStorageVaults"),
		resourceids.UserSpecifiedSegment("exascaleDbStorageVaultName", "exascaleDbStorageVaultName"),
	}
}

// String returns a human-readable description of this Exascale Db Storage Vault ID
func (id ExascaleDbStorageVaultId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Exascale Db Storage Vault Name: %q", id.ExascaleDbStorageVaultName),
	}
	return fmt.Sprintf("Exascale Db Storage Vault (%s)", strings.Join(components, "\n"))
}
