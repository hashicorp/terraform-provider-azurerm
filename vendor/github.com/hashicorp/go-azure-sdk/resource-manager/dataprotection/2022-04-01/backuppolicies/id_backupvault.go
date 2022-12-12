package backuppolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = BackupVaultId{}

// BackupVaultId is a struct representing the Resource ID for a Backup Vault
type BackupVaultId struct {
	SubscriptionId    string
	ResourceGroupName string
	VaultName         string
}

// NewBackupVaultID returns a new BackupVaultId struct
func NewBackupVaultID(subscriptionId string, resourceGroupName string, vaultName string) BackupVaultId {
	return BackupVaultId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VaultName:         vaultName,
	}
}

// ParseBackupVaultID parses 'input' into a BackupVaultId
func ParseBackupVaultID(input string) (*BackupVaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackupVaultId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackupVaultId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, fmt.Errorf("the segment 'vaultName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseBackupVaultIDInsensitively parses 'input' case-insensitively into a BackupVaultId
// note: this method should only be used for API response data and not user input
func ParseBackupVaultIDInsensitively(input string) (*BackupVaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackupVaultId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackupVaultId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, fmt.Errorf("the segment 'vaultName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateBackupVaultID checks that 'input' can be parsed as a Backup Vault ID
func ValidateBackupVaultID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupVaultID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backup Vault ID
func (id BackupVaultId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/backupVaults/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName)
}

// Segments returns a slice of Resource ID Segments which comprise this Backup Vault ID
func (id BackupVaultId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticBackupVaults", "backupVaults", "backupVaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
	}
}

// String returns a human-readable description of this Backup Vault ID
func (id BackupVaultId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
	}
	return fmt.Sprintf("Backup Vault (%s)", strings.Join(components, "\n"))
}
