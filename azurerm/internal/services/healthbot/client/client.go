package client

import (
	"github.com/Azure/azure-sdk-for-go/services/healthbot/mgmt/2020-12-08/healthbot"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	BotClient *healthbot.BotsClient
}

func NewClient(o *common.ClientOptions) *Client {
	botClient := healthbot.NewBotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&botClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		BotClient: &botClient,
	}
}
