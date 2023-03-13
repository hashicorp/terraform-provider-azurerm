package client

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2021-10-01/keyvault" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	keyvaultmgmt "github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type Client struct {
	VaultsClient         *keyvault.VaultsClient
	ManagedHsmClient     *keyvault.ManagedHsmsClient
	ManagementClient     *keyvaultmgmt.BaseClient
	ManagementHSMClient  *keyvaultmgmt.BaseClient
	MHSMSDClient         *keyvaultmgmt.HSMSecurityDomainClient
	MHSMRoleClient       *keyvaultmgmt.RoleDefinitionsClient
	options              *common.ClientOptions
	MHSMRoleAssignClient *keyvaultmgmt.RoleAssignmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	managedHsmClient := keyvault.NewManagedHsmsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedHsmClient.Client, o.ResourceManagerAuthorizer)

	managementClient := keyvaultmgmt.New()
	o.ConfigureClient(&managementClient.Client, o.KeyVaultAuthorizer)

	managementHSMClient := keyvaultmgmt.New()
	o.ConfigureClient(&managementHSMClient.Client, o.ManagedHSMAuthorizer)

	sdClient := keyvaultmgmt.NewHSMSecurityDomainClient()
	o.ConfigureClient(&sdClient.Client, o.ManagedHSMAuthorizer)

	mhsmRoleDefineClient := keyvaultmgmt.NewRoleDefinitionsClient()
	o.ConfigureClient(&mhsmRoleDefineClient.Client, o.ManagedHSMAuthorizer)

	roleAssignClient := keyvaultmgmt.NewRoleAssignmentsClient()
	o.ConfigureClient(&roleAssignClient.Client, o.ManagedHSMAuthorizer)

	vaultsClient := keyvault.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ManagedHsmClient:     &managedHsmClient,
		ManagementClient:     &managementClient,
		ManagementHSMClient:  &managementHSMClient,
		MHSMSDClient:         &sdClient,
		VaultsClient:         &vaultsClient,
		MHSMRoleClient:       &mhsmRoleDefineClient,
		MHSMRoleAssignClient: &roleAssignClient,
		options:              o,
	}
}

func (c *Client) KeyVaultClientForSubscription(subscriptionId string) *keyvault.VaultsClient {
	vaultsClient := keyvault.NewVaultsClientWithBaseURI(c.options.ResourceManagerEndpoint, subscriptionId)
	c.options.ConfigureClient(&vaultsClient.Client, c.options.ResourceManagerAuthorizer)
	return &vaultsClient
}
