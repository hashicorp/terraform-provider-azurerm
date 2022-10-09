package client

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2021-10-01/keyvault" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	keyvaultmgmt "github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type Client struct {
	ManagedHsmClient *keyvault.ManagedHsmsClient
	ManagementClient *keyvaultmgmt.BaseClient
	MHSMSDClient     *keyvaultmgmt.HSMSecurityDomainClient
	VaultsClient     *keyvault.VaultsClient
	MHSMRoleClient   *keyvaultmgmt.RoleDefinitionsClient
	options          *common.ClientOptions
}

func NewClient(o *common.ClientOptions) *Client {
	managedHsmClient := keyvault.NewManagedHsmsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedHsmClient.Client, o.ResourceManagerAuthorizer)

	managementClient := keyvaultmgmt.New()
	o.ConfigureClient(&managementClient.Client, o.KeyVaultAuthorizer)

	sdClient := keyvaultmgmt.NewHSMSecurityDomainClient()
	o.ConfigureClient(&sdClient.Client, o.ManagedHSMAuthorizer)

	mhsmRoleDefineClient := keyvaultmgmt.NewRoleDefinitionsClient()
	o.ConfigureClient(&mhsmRoleDefineClient.Client, o.ManagedHSMAuthorizer)

	vaultsClient := keyvault.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ManagedHsmClient: &managedHsmClient,
		ManagementClient: &managementClient,
		MHSMSDClient:     &sdClient,
		VaultsClient:     &vaultsClient,
		MHSMRoleClient:   &mhsmRoleDefineClient,
		options:          o,
	}
}

func (client Client) KeyVaultClientForSubscription(subscriptionId string) *keyvault.VaultsClient {
	vaultsClient := keyvault.NewVaultsClientWithBaseURI(client.options.ResourceManagerEndpoint, subscriptionId)
	client.options.ConfigureClient(&vaultsClient.Client, client.options.ResourceManagerAuthorizer)
	return &vaultsClient
}
