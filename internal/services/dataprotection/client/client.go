package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/resourceguards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/legacysdk/dataprotection"
)

type Client struct {
	BackupVaultClient    *dataprotection.BackupVaultsClient
	BackupPolicyClient   *dataprotection.BackupPoliciesClient
	BackupInstanceClient *dataprotection.BackupInstancesClient
	ResourceGuardClient  *resourceguards.ResourceGuardsClient
}

func NewClient(o *common.ClientOptions) *Client {
	backupVaultClient := dataprotection.NewBackupVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backupVaultClient.Client, o.ResourceManagerAuthorizer)

	backupPolicyClient := dataprotection.NewBackupPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backupPolicyClient.Client, o.ResourceManagerAuthorizer)

	backupInstanceClient := dataprotection.NewBackupInstancesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
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
