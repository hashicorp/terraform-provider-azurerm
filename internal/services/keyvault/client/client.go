package client

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2021-10-01/keyvault" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2021-10-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	dataplane "github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type Client struct {
	ManagedHsmClient *keyvault.ManagedHsmsClient
	ManagementClient *dataplane.BaseClient
	VaultsClient     *vaults.VaultsClient
}

func NewClient(o *common.ClientOptions) *Client {
	managedHsmClient := keyvault.NewManagedHsmsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedHsmClient.Client, o.ResourceManagerAuthorizer)

	managementClient := dataplane.New()
	o.ConfigureClient(&managementClient.Client, o.KeyVaultAuthorizer)

	vaultsClient := vaults.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ManagedHsmClient: &managedHsmClient,
		ManagementClient: &managementClient,
		VaultsClient:     &vaultsClient,
	}
}
