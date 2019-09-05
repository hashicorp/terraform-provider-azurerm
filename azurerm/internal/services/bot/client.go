package bot

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	BotClient *botservice.BotsClient
}

func BuildClient(o *common.ClientOptions) *Client {

	BotClient := botservice.NewBotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&BotClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		BotClient: &BotClient,
	}
}
