// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &KeyVaultKeyId{}

// KeyVaultKeyId is a struct representing the Resource ID for a Key
type KeyVaultKeyId struct {
	SubscriptionId    string
	ResourceGroupName string
	VaultName         string
	KeyName           string
}

// NewKeyVaultKeyID returns a new KeyVaultKeyId struct
func NewKeyVaultKeyID(subscriptionId string, resourceGroupName string, vaultName string, keyName string) KeyVaultKeyId {
	return KeyVaultKeyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VaultName:         vaultName,
		KeyName:           keyName,
	}
}

// ParseKeyVaultKeyID parses 'input' into a KeyVaultKeyId
func ParseKeyVaultKeyID(input string) (*KeyVaultKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KeyVaultKeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KeyVaultKeyId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseKeyVaultKeyIDInsensitively parses 'input' case-insensitively into a KeyVaultKeyId
// note: this method should only be used for API response data and not user input
func ParseKeyVaultKeyIDInsensitively(input string) (*KeyVaultKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KeyVaultKeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KeyVaultKeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *KeyVaultKeyId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.KeyName, ok = input.Parsed["keyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "keyName", input)
	}

	return nil
}

// ValidateKeyVaultKeyID checks that 'input' can be parsed as a Key ID
func ValidateKeyVaultKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKeyVaultKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Key ID
func (id KeyVaultKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/keys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.KeyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Key ID
func (id KeyVaultKeyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKeyVault", "Microsoft.KeyVault", "Microsoft.KeyVault"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
		resourceids.StaticSegment("staticKeys", "keys", "keys"),
		resourceids.UserSpecifiedSegment("keyName", "keyValue"),
	}
}

// String returns a human-readable description of this Key ID
func (id KeyVaultKeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Key Vault Name: %q", id.VaultName),
		fmt.Sprintf("Key Name: %q", id.KeyName),
	}
	return fmt.Sprintf("Key Vault Key (%s)", strings.Join(components, "\n"))
}
