package recoveryservices

import (
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2017-07-01/backup"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ProtectedItemsClient     *backup.ProtectedItemsGroupClient
	ProtectionPoliciesClient *backup.ProtectionPoliciesClient
	VaultsClient             *recoveryservices.VaultsClient
}

func BuildClient(o *common.ClientOptions) *Client {

	VaultsClient := recoveryservices.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VaultsClient.Client, o.ResourceManagerAuthorizer)

	ProtectedItemsClient := backup.NewProtectedItemsGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ProtectedItemsClient.Client, o.ResourceManagerAuthorizer)

	ProtectionPoliciesClient := backup.NewProtectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ProtectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ProtectedItemsClient:     &ProtectedItemsClient,
		ProtectionPoliciesClient: &ProtectionPoliciesClient,
		VaultsClient:             &VaultsClient,
	}
}
