// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ProtectionContainerId struct {
	SubscriptionId   string
	ResourceGroup    string
	VaultName        string
	BackupFabricName string
	Name             string
}

func NewProtectionContainerID(subscriptionId, resourceGroup, vaultName, backupFabricName, name string) ProtectionContainerId {
	return ProtectionContainerId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		VaultName:        vaultName,
		BackupFabricName: backupFabricName,
		Name:             name,
	}
}

func (id ProtectionContainerId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Backup Fabric Name %q", id.BackupFabricName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Protection Container", segmentsStr)
}

func (id ProtectionContainerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/backupFabrics/%s/protectionContainers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.BackupFabricName, id.Name)
}

// ProtectionContainerID parses a ProtectionContainer ID into an ProtectionContainerId struct
func ProtectionContainerID(input string) (*ProtectionContainerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ProtectionContainer ID: %+v", input, err)
	}

	resourceId := ProtectionContainerId{
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
	if resourceId.Name, err = id.PopSegment("protectionContainers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
