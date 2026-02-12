// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/resourceguardresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BackupVaultClient    *backupvaultresources.BackupVaultResourcesClient
	BackupPolicyClient   *basebackuppolicyresources.BaseBackupPolicyResourcesClient
	BackupInstanceClient *backupinstanceresources.BackupInstanceResourcesClient
	ResourceGuardClient  *resourceguardresources.ResourceGuardResourcesClient
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

	return &Client{
		BackupVaultClient:    backupVaultClient,
		BackupPolicyClient:   backupPolicyClient,
		BackupInstanceClient: backupInstanceClient,
		ResourceGuardClient:  resourceGuardClient,
	}, nil
}
