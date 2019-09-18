package analysisservices

import (
	"github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-08-01/analysisservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ServerClient *analysisservices.ServersClient
}

func BuildClient(o *common.ClientOptions) *Client {

	ServerClient := analysisservices.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServerClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServerClient: &ServerClient,
	}
}
