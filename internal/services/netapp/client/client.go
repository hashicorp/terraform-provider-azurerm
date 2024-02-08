// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/capacitypools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/netappaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/snapshotpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/volumequotarules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/volumes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/volumesreplication"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountClient           *netappaccounts.NetAppAccountsClient
	PoolClient              *capacitypools.CapacityPoolsClient
	VolumeClient            *volumes.VolumesClient
	VolumeGroupClient       *volumegroups.VolumeGroupsClient
	VolumeReplicationClient *volumesreplication.VolumesReplicationClient
	VolumeQuotaRules        *volumequotarules.VolumeQuotaRulesClient
	SnapshotClient          *snapshots.SnapshotsClient
	SnapshotPoliciesClient  *snapshotpolicy.SnapshotPolicyClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	accountClient, err := netappaccounts.NewNetAppAccountsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(accountClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AutomationAccount client: %+v", err)
	}

	poolClient, err := capacitypools.NewCapacityPoolsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(poolClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PoolClient client: %+v", err)
	}

	volumeClient, err := volumes.NewVolumesClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(volumeClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VolumeClient client: %+v", err)
	}

	volumeGroupClient, err := volumegroups.NewVolumeGroupsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(volumeGroupClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VolumeGroupClient client: %+v", err)
	}

	volumeReplicationClient, err := volumesreplication.NewVolumesReplicationClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(volumeReplicationClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VolumeReplicationClient client: %+v", err)
	}

	volumeQuotaRuleClient, err := volumequotarules.NewVolumeQuotaRulesClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(volumeQuotaRuleClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VolumeQuotaRuleClient client: %+v", err)
	}

	snapshotClient, err := snapshots.NewSnapshotsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(snapshotClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SnapshotClient client: %+v", err)
	}

	snapshotPoliciesClient, err := snapshotpolicy.NewSnapshotPolicyClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(snapshotPoliciesClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SnapshotPoliciesClient client: %+v", err)
	}

	return &Client{
		AccountClient:           accountClient,
		PoolClient:              poolClient,
		VolumeClient:            volumeClient,
		VolumeGroupClient:       volumeGroupClient,
		VolumeReplicationClient: volumeReplicationClient,
		VolumeQuotaRules:        volumeQuotaRuleClient,
		SnapshotClient:          snapshotClient,
		SnapshotPoliciesClient:  snapshotPoliciesClient,
	}, nil
}
