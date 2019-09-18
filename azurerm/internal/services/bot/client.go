package bot

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	BotClient        *botservice.BotsClient
	ConnectionClient *botservice.BotConnectionClient
}

func BuildClient(o *common.ClientOptions) *Client {

	botClient := botservice.NewBotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&botClient.Client, o.ResourceManagerAuthorizer)

	connectionClient := botservice.NewBotConnectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&connectionClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		BotClient:        &botClient,
		ConnectionClient: &connectionClient,
	}
}
