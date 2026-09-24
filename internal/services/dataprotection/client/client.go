// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources"
	backupinstanceresources20260601 "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	backupinstances20260601 "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstances"
	basebackuppolicyresources20260601 "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BackupVaultClient    *backupvaultresources.BackupVaultResourcesClient
	BackupPolicyClient   *basebackuppolicyresources.BaseBackupPolicyResourcesClient
	BackupInstanceClient *backupinstanceresources.BackupInstanceResourcesClient
	ResourceGuardClient  *resourceguardresources.ResourceGuardResourcesClient

	BackupPolicyClient20260601    *basebackuppolicyresources20260601.BaseBackupPolicyResourcesClient
	BackupInstanceClient20260601  *backupinstanceresources20260601.BackupInstanceResourcesClient
	BackupInstancesClient20260601 *backupinstances20260601.BackupInstancesClient
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

	resourceGuardClient, err := resourceguardresources.NewResourceGuardResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ResourceGuardResources client: %+v", err)
	}
	o.Configure(resourceGuardClient.Client, o.Authorizers.ResourceManager)

	backupPolicyClient20260601, err := basebackuppolicyresources20260601.NewBaseBackupPolicyResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building 2026-06-01 BaseBackupPolicyResources client: %+v", err)
	}
	o.Configure(backupPolicyClient20260601.Client, o.Authorizers.ResourceManager)

	backupInstanceClient20260601, err := backupinstanceresources20260601.NewBackupInstanceResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building 2026-06-01 BackupInstanceResources client: %+v", err)
	}
	o.Configure(backupInstanceClient20260601.Client, o.Authorizers.ResourceManager)

	backupInstancesClient20260601, err := backupinstances20260601.NewBackupInstancesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building 2026-06-01 BackupInstances client: %+v", err)
	}
	o.Configure(backupInstancesClient20260601.Client, o.Authorizers.ResourceManager)

	return &Client{
		BackupVaultClient:    backupVaultClient,
		BackupPolicyClient:   backupPolicyClient,
		BackupInstanceClient: backupInstanceClient,
		ResourceGuardClient:  resourceGuardClient,

		BackupPolicyClient20260601:    backupPolicyClient20260601,
		BackupInstanceClient20260601:  backupInstanceClient20260601,
		BackupInstancesClient20260601: backupInstancesClient20260601,
	}, nil
}
