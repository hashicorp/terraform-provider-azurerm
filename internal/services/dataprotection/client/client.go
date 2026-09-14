// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources"
	backupinstanceresources20260301Client "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-03-01/backupinstanceresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BackupVaultClient    *backupvaultresources.BackupVaultResourcesClient
	BackupPolicyClient   *basebackuppolicyresources.BaseBackupPolicyResourcesClient
	BackupInstanceClient *backupinstanceresources.BackupInstanceResourcesClient
	ResourceGuardClient  *resourceguardresources.ResourceGuardResourcesClient

	// BackupInstanceClient_v2026_03_01 is used by `azurerm_data_protection_backup_instance_blob_storage` only: 2026-03-01
	// is the first API version with the blob auto protection models. The rest of the service package deliberately stays
	// on 2025-07-01 since newer API versions enforce `AlwaysOn` soft delete on the backup vault, see #31874 / #31877.
	BackupInstanceClient_v2026_03_01 *backupinstanceresources20260301Client.BackupInstanceResourcesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	backupVaultClient, err := backupvaultresources.NewBackupVaultResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BackupVaultResources client: %+v", err)
	}
	o.Configure(backupVaultClient.Client, o.Authorizers.ResourceManager)

	backupPolicyClient, err := basebackuppolicyresources.NewBaseBackupPolicyResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BaseBackupPolicyResources client: %+v", err)
	}
	o.Configure(backupPolicyClient.Client, o.Authorizers.ResourceManager)

	backupInstanceClient, err := backupinstanceresources.NewBackupInstanceResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BackupInstances client: %+v", err)
	}
	o.Configure(backupInstanceClient.Client, o.Authorizers.ResourceManager)

	backupInstanceClient_v2026_03_01, err := backupinstanceresources20260301Client.NewBackupInstanceResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BackupInstances (2026-03-01) client: %+v", err)
	}
	o.Configure(backupInstanceClient_v2026_03_01.Client, o.Authorizers.ResourceManager)

	resourceGuardClient, err := resourceguardresources.NewResourceGuardResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ResourceGuardResources client: %+v", err)
	}
	o.Configure(resourceGuardClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		BackupVaultClient:    backupVaultClient,
		BackupPolicyClient:   backupPolicyClient,
		BackupInstanceClient: backupInstanceClient,
		ResourceGuardClient:  resourceGuardClient,

		BackupInstanceClient_v2026_03_01: backupInstanceClient_v2026_03_01,
	}, nil
}
