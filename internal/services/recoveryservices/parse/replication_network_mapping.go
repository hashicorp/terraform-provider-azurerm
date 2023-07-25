// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ReplicationNetworkMappingId struct {
	SubscriptionId         string
	ResourceGroup          string
	VaultName              string
	ReplicationFabricName  string
	ReplicationNetworkName string
	Name                   string
}

func NewReplicationNetworkMappingID(subscriptionId, resourceGroup, vaultName, replicationFabricName, replicationNetworkName, name string) ReplicationNetworkMappingId {
	return ReplicationNetworkMappingId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		VaultName:              vaultName,
		ReplicationFabricName:  replicationFabricName,
		ReplicationNetworkName: replicationNetworkName,
		Name:                   name,
	}
}

func (id ReplicationNetworkMappingId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Replication Network Name %q", id.ReplicationNetworkName),
		fmt.Sprintf("Replication Fabric Name %q", id.ReplicationFabricName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Replication Network Mapping", segmentsStr)
}

func (id ReplicationNetworkMappingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s/replicationNetworks/%s/replicationNetworkMappings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.ReplicationFabricName, id.ReplicationNetworkName, id.Name)
}

// ReplicationNetworkMappingID parses a ReplicationNetworkMapping ID into an ReplicationNetworkMappingId struct
func ReplicationNetworkMappingID(input string) (*ReplicationNetworkMappingId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ReplicationNetworkMapping ID: %+v", input, err)
	}

	resourceId := ReplicationNetworkMappingId{
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
	if resourceId.ReplicationNetworkName, err = id.PopSegment("replicationNetworks"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("replicationNetworkMappings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
