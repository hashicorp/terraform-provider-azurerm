package client

import (
	keyvaultmgmt "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ManagementClient *keyvaultmgmt.BaseClient
	VaultsClient     *keyvault.VaultsClient
}

func NewClient(o *common.ClientOptions) *Client {
	managementClient := keyvaultmgmt.New()
	o.ConfigureClient(&managementClient.Client, o.KeyVaultAuthorizer)

	vaultsClient := keyvault.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ManagementClient: &managementClient,
		VaultsClient:     &vaultsClient,
	}
}
