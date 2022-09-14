package client

import (
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

func NewClient(o *common.ClientOptions) *Client {
	backupVaultClient := backupvaults.NewBackupVaultsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&backupVaultClient.Client, o.ResourceManagerAuthorizer)

	backupPolicyClient := backuppolicies.NewBackupPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&backupPolicyClient.Client, o.ResourceManagerAuthorizer)

	backupInstanceClient := backupinstances.NewBackupInstancesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&backupInstanceClient.Client, o.ResourceManagerAuthorizer)

	resourceGuardClient := resourceguards.NewResourceGuardsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&resourceGuardClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		BackupVaultClient:    &backupVaultClient,
		BackupPolicyClient:   &backupPolicyClient,
		BackupInstanceClient: &backupInstanceClient,
		ResourceGuardClient:  &resourceGuardClient,
	}
}
