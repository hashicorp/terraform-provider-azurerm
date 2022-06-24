package operationstatus

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = BackupVaultOperationStatuId{}

// BackupVaultOperationStatuId is a struct representing the Resource ID for a Backup Vault Operation Statu
type BackupVaultOperationStatuId struct {
	SubscriptionId    string
	ResourceGroupName string
	VaultName         string
	OperationId       string
}

// NewBackupVaultOperationStatuID returns a new BackupVaultOperationStatuId struct
func NewBackupVaultOperationStatuID(subscriptionId string, resourceGroupName string, vaultName string, operationId string) BackupVaultOperationStatuId {
	return BackupVaultOperationStatuId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VaultName:         vaultName,
		OperationId:       operationId,
	}
}

// ParseBackupVaultOperationStatuID parses 'input' into a BackupVaultOperationStatuId
func ParseBackupVaultOperationStatuID(input string) (*BackupVaultOperationStatuId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackupVaultOperationStatuId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackupVaultOperationStatuId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, fmt.Errorf("the segment 'vaultName' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseBackupVaultOperationStatuIDInsensitively parses 'input' case-insensitively into a BackupVaultOperationStatuId
// note: this method should only be used for API response data and not user input
func ParseBackupVaultOperationStatuIDInsensitively(input string) (*BackupVaultOperationStatuId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackupVaultOperationStatuId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackupVaultOperationStatuId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, fmt.Errorf("the segment 'vaultName' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateBackupVaultOperationStatuID checks that 'input' can be parsed as a Backup Vault Operation Statu ID
func ValidateBackupVaultOperationStatuID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupVaultOperationStatuID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backup Vault Operation Statu ID
func (id BackupVaultOperationStatuId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/backupVaults/%s/operationStatus/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Backup Vault Operation Statu ID
func (id BackupVaultOperationStatuId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticBackupVaults", "backupVaults", "backupVaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
		resourceids.StaticSegment("staticOperationStatus", "operationStatus", "operationStatus"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Backup Vault Operation Statu ID
func (id BackupVaultOperationStatuId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Backup Vault Operation Statu (%s)", strings.Join(components, "\n"))
}
