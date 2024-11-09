// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ReplicationProtectionContainerId struct {
	SubscriptionId        string
	ResourceGroup         string
	VaultName             string
	ReplicationFabricName string
	Name                  string
}

func NewReplicationProtectionContainerID(subscriptionId, resourceGroup, vaultName, replicationFabricName, name string) ReplicationProtectionContainerId {
	return ReplicationProtectionContainerId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		VaultName:             vaultName,
		ReplicationFabricName: replicationFabricName,
		Name:                  name,
	}
}

func (id ReplicationProtectionContainerId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Replication Fabric Name %q", id.ReplicationFabricName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Replication Protection Container", segmentsStr)
}

func (id ReplicationProtectionContainerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s/replicationProtectionContainers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.ReplicationFabricName, id.Name)
}

// ReplicationProtectionContainerID parses a ReplicationProtectionContainer ID into an ReplicationProtectionContainerId struct
func ReplicationProtectionContainerID(input string) (*ReplicationProtectionContainerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ReplicationProtectionContainer ID: %+v", input, err)
	}

	resourceId := ReplicationProtectionContainerId{
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
	if resourceId.Name, err = id.PopSegment("replicationProtectionContainers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
