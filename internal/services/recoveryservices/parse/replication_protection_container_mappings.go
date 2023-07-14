// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ReplicationProtectionContainerMappingsId struct {
	SubscriptionId                            string
	ResourceGroup                             string
	VaultName                                 string
	ReplicationFabricName                     string
	ReplicationProtectionContainerName        string
	ReplicationProtectionContainerMappingName string
}

func NewReplicationProtectionContainerMappingsID(subscriptionId, resourceGroup, vaultName, replicationFabricName, replicationProtectionContainerName, replicationProtectionContainerMappingName string) ReplicationProtectionContainerMappingsId {
	return ReplicationProtectionContainerMappingsId{
		SubscriptionId:                     subscriptionId,
		ResourceGroup:                      resourceGroup,
		VaultName:                          vaultName,
		ReplicationFabricName:              replicationFabricName,
		ReplicationProtectionContainerName: replicationProtectionContainerName,
		ReplicationProtectionContainerMappingName: replicationProtectionContainerMappingName,
	}
}

func (id ReplicationProtectionContainerMappingsId) String() string {
	segments := []string{
		fmt.Sprintf("Replication Protection Container Mapping Name %q", id.ReplicationProtectionContainerMappingName),
		fmt.Sprintf("Replication Protection Container Name %q", id.ReplicationProtectionContainerName),
		fmt.Sprintf("Replication Fabric Name %q", id.ReplicationFabricName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Replication Protection Container Mappings", segmentsStr)
}

func (id ReplicationProtectionContainerMappingsId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s/replicationProtectionContainers/%s/replicationProtectionContainerMappings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.ReplicationFabricName, id.ReplicationProtectionContainerName, id.ReplicationProtectionContainerMappingName)
}

// ReplicationProtectionContainerMappingsID parses a ReplicationProtectionContainerMappings ID into an ReplicationProtectionContainerMappingsId struct
func ReplicationProtectionContainerMappingsID(input string) (*ReplicationProtectionContainerMappingsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ReplicationProtectionContainerMappings ID: %+v", input, err)
	}

	resourceId := ReplicationProtectionContainerMappingsId{
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
	if resourceId.ReplicationProtectionContainerMappingName, err = id.PopSegment("replicationProtectionContainerMappings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
