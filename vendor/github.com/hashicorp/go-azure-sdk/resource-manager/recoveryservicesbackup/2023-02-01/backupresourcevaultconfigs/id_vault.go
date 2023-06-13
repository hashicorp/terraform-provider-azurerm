package backupresourcevaultconfigs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VaultId{}

// VaultId is a struct representing the Resource ID for a Vault
type VaultId struct {
	SubscriptionId    string
	ResourceGroupName string
	VaultName         string
}

// NewVaultID returns a new VaultId struct
func NewVaultID(subscriptionId string, resourceGroupName string, vaultName string) VaultId {
	return VaultId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VaultName:         vaultName,
	}
}

// ParseVaultID parses 'input' into a VaultId
func ParseVaultID(input string) (*VaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(VaultId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VaultId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	return &id, nil
}

// ParseVaultIDInsensitively parses 'input' case-insensitively into a VaultId
// note: this method should only be used for API response data and not user input
func ParseVaultIDInsensitively(input string) (*VaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(VaultId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VaultId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	return &id, nil
}

// ValidateVaultID checks that 'input' can be parsed as a Vault ID
func ValidateVaultID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVaultID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Vault ID
func (id VaultId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName)
}

// Segments returns a slice of Resource ID Segments which comprise this Vault ID
func (id VaultId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRecoveryServices", "Microsoft.RecoveryServices", "Microsoft.RecoveryServices"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
	}
}

// String returns a human-readable description of this Vault ID
func (id VaultId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
	}
	return fmt.Sprintf("Vault (%s)", strings.Join(components, "\n"))
}
