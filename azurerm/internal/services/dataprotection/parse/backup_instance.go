package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BackupInstanceId struct {
	SubscriptionId  string
	ResourceGroup   string
	BackupVaultName string
	Name            string
}

func NewBackupInstanceID(subscriptionId, resourceGroup, backupVaultName, name string) BackupInstanceId {
	return BackupInstanceId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		BackupVaultName: backupVaultName,
		Name:            name,
	}
}

func (id BackupInstanceId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Backup Vault Name %q", id.BackupVaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Backup Instance", segmentsStr)
}

func (id BackupInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/backupVaults/%s/backupInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.BackupVaultName, id.Name)
}

// BackupInstanceID parses a BackupInstance ID into an BackupInstanceId struct
func BackupInstanceID(input string) (*BackupInstanceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := BackupInstanceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.BackupVaultName, err = id.PopSegment("backupVaults"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("backupInstances"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
