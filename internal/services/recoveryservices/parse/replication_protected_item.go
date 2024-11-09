// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ReplicationProtectedItemId struct {
	SubscriptionId                     string
	ResourceGroup                      string
	VaultName                          string
	ReplicationFabricName              string
	ReplicationProtectionContainerName string
	Name                               string
}

func NewReplicationProtectedItemID(subscriptionId, resourceGroup, vaultName, replicationFabricName, replicationProtectionContainerName, name string) ReplicationProtectedItemId {
	return ReplicationProtectedItemId{
		SubscriptionId:                     subscriptionId,
		ResourceGroup:                      resourceGroup,
		VaultName:                          vaultName,
		ReplicationFabricName:              replicationFabricName,
		ReplicationProtectionContainerName: replicationProtectionContainerName,
		Name:                               name,
	}
}

func (id ReplicationProtectedItemId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Replication Protection Container Name %q", id.ReplicationProtectionContainerName),
		fmt.Sprintf("Replication Fabric Name %q", id.ReplicationFabricName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Replication Protected Item", segmentsStr)
}

func (id ReplicationProtectedItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s/replicationProtectionContainers/%s/replicationProtectedItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.ReplicationFabricName, id.ReplicationProtectionContainerName, id.Name)
}

// ReplicationProtectedItemID parses a ReplicationProtectedItem ID into an ReplicationProtectedItemId struct
func ReplicationProtectedItemID(input string) (*ReplicationProtectedItemId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ReplicationProtectedItem ID: %+v", input, err)
	}

	resourceId := ReplicationProtectedItemId{
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
	if resourceId.ReplicationFabricName, err = id.PopSegment("replicationFabrics"); err != nil {
		return nil, err
	}
	if resourceId.ReplicationProtectionContainerName, err = id.PopSegment("replicationProtectionContainers"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("replicationProtectedItems"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
