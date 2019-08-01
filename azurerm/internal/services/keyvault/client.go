package keyvault

import (
	keyvaultmgmt "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	VaultsClient     keyvault.VaultsClient
	ManagementClient keyvaultmgmt.BaseClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.VaultsClient = keyvault.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VaultsClient.Client, o.ResourceManagerAuthorizer)

	c.ManagementClient = keyvaultmgmt.New()
	o.ConfigureClient(&c.ManagementClient.Client, o.KeyVaultAuthorizer)

	return &c
}
