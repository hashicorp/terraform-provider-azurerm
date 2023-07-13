// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backupinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backupvaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/resourceguards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BackupVaultClient    *backupvaults.BackupVaultsClient
	BackupPolicyClient   *backuppolicies.BackupPoliciesClient
	BackupInstanceClient *backupinstances.BackupInstancesClient
	ResourceGuardClient  *resourceguards.ResourceGuardsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	backupVaultClient, err := backupvaults.NewBackupVaultsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BackupVaults client: %+v", err)
	}
	o.Configure(backupVaultClient.Client, o.Authorizers.ResourceManager)

	backupPolicyClient, err := backuppolicies.NewBackupPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BackupPolicies client: %+v", err)
	}
	o.Configure(backupPolicyClient.Client, o.Authorizers.ResourceManager)

	backupInstanceClient, err := backupinstances.NewBackupInstancesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BackupInstances client: %+v", err)
	}
	o.Configure(backupInstanceClient.Client, o.Authorizers.ResourceManager)

	resourceGuardClient, err := resourceguards.NewResourceGuardsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ResourceGuards client: %+v", err)
	}
	o.Configure(resourceGuardClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		BackupVaultClient:    backupVaultClient,
		BackupPolicyClient:   backupPolicyClient,
		BackupInstanceClient: backupInstanceClient,
		ResourceGuardClient:  resourceGuardClient,
	}, nil
}
