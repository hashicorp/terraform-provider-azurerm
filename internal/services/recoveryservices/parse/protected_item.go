// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ProtectedItemId struct {
	SubscriptionId          string
	ResourceGroup           string
	VaultName               string
	BackupFabricName        string
	ProtectionContainerName string
	Name                    string
}

func NewProtectedItemID(subscriptionId, resourceGroup, vaultName, backupFabricName, protectionContainerName, name string) ProtectedItemId {
	return ProtectedItemId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		VaultName:               vaultName,
		BackupFabricName:        backupFabricName,
		ProtectionContainerName: protectionContainerName,
		Name:                    name,
	}
}

func (id ProtectedItemId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Protection Container Name %q", id.ProtectionContainerName),
		fmt.Sprintf("Backup Fabric Name %q", id.BackupFabricName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Protected Item", segmentsStr)
}

func (id ProtectedItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/backupFabrics/%s/protectionContainers/%s/protectedItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.BackupFabricName, id.ProtectionContainerName, id.Name)
}

// ProtectedItemID parses a ProtectedItem ID into an ProtectedItemId struct
func ProtectedItemID(input string) (*ProtectedItemId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ProtectedItem ID: %+v", input, err)
	}

	resourceId := ProtectedItemId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VaultName, err = id.PopSegment("vaults"); err != nil {
		return nil, err
	}
	if resourceId.BackupFabricName, err = id.PopSegment("backupFabrics"); err != nil {
		return nil, err
	}
	if resourceId.ProtectionContainerName, err = id.PopSegment("protectionContainers"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("protectedItems"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
