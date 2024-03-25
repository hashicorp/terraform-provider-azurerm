// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &KeyVaultId{}

// KeyVaultId is a struct representing the Resource ID for a Key Vault
type KeyVaultId struct {
	SubscriptionId    string
	ResourceGroupName string
	VaultName         string
}

// NewKeyVaultID returns a new KeyVaultId struct
func NewKeyVaultID(subscriptionId string, resourceGroupName string, vaultName string) KeyVaultId {
	return KeyVaultId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VaultName:         vaultName,
	}
}

// ParseKeyVaultID parses 'input' into a KeyVaultId
func ParseKeyVaultID(input string) (*KeyVaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KeyVaultId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KeyVaultId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseKeyVaultIDInsensitively parses 'input' case-insensitively into a KeyVaultId
// note: this method should only be used for API response data and not user input
func ParseKeyVaultIDInsensitively(input string) (*KeyVaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KeyVaultId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KeyVaultId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *KeyVaultId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VaultName, ok = input.Parsed["vaultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vaultName", input)
	}

	return nil
}

// ValidateKeyVaultID checks that 'input' can be parsed as a Key Vault ID
func ValidateKeyVaultID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKeyVaultID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Key Vault ID
func (id KeyVaultId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName)
}

// Segments returns a slice of Resource ID Segments which comprise this Key Vault ID
func (id KeyVaultId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKeyVault", "Microsoft.KeyVault", "Microsoft.KeyVault"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
	}
}

// String returns a human-readable description of this Key Vault ID
func (id KeyVaultId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Key Vault Name: %q", id.VaultName),
	}
	return fmt.Sprintf("Key Vault (%s)", strings.Join(components, "\n"))
}
