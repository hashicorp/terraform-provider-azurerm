// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/backuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/backups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/backupvaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/capacitypools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/netappaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/snapshotpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/volumequotarules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountClient          *netappaccounts.NetAppAccountsClient
	PoolClient             *capacitypools.CapacityPoolsClient
	VolumeClient           *volumes.VolumesClient
	VolumeGroupClient      *volumegroups.VolumeGroupsClient
	VolumeQuotaRules       *volumequotarules.VolumeQuotaRulesClient
	SnapshotClient         *snapshots.SnapshotsClient
	SnapshotPoliciesClient *snapshotpolicies.SnapshotPoliciesClient
	BackupVaultsClient     *backupvaults.BackupVaultsClient
	BackupPolicyClient     *backuppolicies.BackupPoliciesClient
	BackupClient           *backups.BackupsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	accountClient, err := netappaccounts.NewNetAppAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AutomationAccount client: %+v", err)
	}
	o.Configure(accountClient.Client, o.Authorizers.ResourceManager)

	poolClient, err := capacitypools.NewCapacityPoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PoolClient client: %+v", err)
	}
	o.Configure(poolClient.Client, o.Authorizers.ResourceManager)

	volumeClient, err := volumes.NewVolumesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VolumeClient client: %+v", err)
	}
	o.Configure(volumeClient.Client, o.Authorizers.ResourceManager)

	volumeGroupClient, err := volumegroups.NewVolumeGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VolumeGroupClient client: %+v", err)
	}
	o.Configure(volumeGroupClient.Client, o.Authorizers.ResourceManager)

	volumeQuotaRuleClient, err := volumequotarules.NewVolumeQuotaRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VolumeQuotaRuleClient client: %+v", err)
	}
	o.Configure(volumeQuotaRuleClient.Client, o.Authorizers.ResourceManager)

	snapshotClient, err := snapshots.NewSnapshotsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SnapshotClient client: %+v", err)
	}
	o.Configure(snapshotClient.Client, o.Authorizers.ResourceManager)

	snapshotPoliciesClient, err := snapshotpolicies.NewSnapshotPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SnapshotPoliciesClient client: %+v", err)
	}
	o.Configure(snapshotPoliciesClient.Client, o.Authorizers.ResourceManager)

	backupVaultsClient, err := backupvaults.NewBackupVaultsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BackupVaultsClient client: %+v", err)
	}
	o.Configure(backupVaultsClient.Client, o.Authorizers.ResourceManager)

	backupPolicyClient, err := backuppolicies.NewBackupPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BackupPoliciesClient client: %+v", err)
	}
	o.Configure(backupPolicyClient.Client, o.Authorizers.ResourceManager)

	backupClient, err := backups.NewBackupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BackupClient client: %+v", err)
	}
	o.Configure(backupClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AccountClient:          accountClient,
		PoolClient:             poolClient,
		VolumeClient:           volumeClient,
		VolumeGroupClient:      volumeGroupClient,
		VolumeQuotaRules:       volumeQuotaRuleClient,
		SnapshotClient:         snapshotClient,
		SnapshotPoliciesClient: snapshotPoliciesClient,
		BackupVaultsClient:     backupVaultsClient,
		BackupPolicyClient:     backupPolicyClient,
		BackupClient:           backupClient,
	}, nil
}
