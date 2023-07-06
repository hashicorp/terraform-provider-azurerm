// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ReplicationFabricId struct {
	SubscriptionId string
	ResourceGroup  string
	VaultName      string
	Name           string
}

func NewReplicationFabricID(subscriptionId, resourceGroup, vaultName, name string) ReplicationFabricId {
	return ReplicationFabricId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VaultName:      vaultName,
		Name:           name,
	}
}

func (id ReplicationFabricId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Replication Fabric", segmentsStr)
}

func (id ReplicationFabricId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.Name)
}

// ReplicationFabricID parses a ReplicationFabric ID into an ReplicationFabricId struct
func ReplicationFabricID(input string) (*ReplicationFabricId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ReplicationFabric ID: %+v", input, err)
	}

	resourceId := ReplicationFabricId{
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
	if resourceId.Name, err = id.PopSegment("replicationFabrics"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
