package client

import (
	"github.com/Azure/azure-sdk-for-go/services/msi/mgmt/2018-11-30/msi"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	UserAssignedIdentitiesClient *msi.UserAssignedIdentitiesClient
}

func NewClient(o *common.ClientOptions) *Client {
	UserAssignedIdentitiesClient := msi.NewUserAssignedIdentitiesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&UserAssignedIdentitiesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		UserAssignedIdentitiesClient: &UserAssignedIdentitiesClient,
	}
}
