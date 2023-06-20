package client

import (
	"github.com/Azure/go-autorest/autorest"
	healthbot_2022_08_08 "github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2022-08-08"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

type Client struct {
	BotClient        *botservice.BotsClient
	ConnectionClient *botservice.BotConnectionClient
	ChannelClient    *botservice.ChannelsClient
	HealthBotClient  *healthbot_2022_08_08.Client
}

func NewClient(o *common.ClientOptions) *Client {
	botClient := botservice.NewBotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&botClient.Client, o.ResourceManagerAuthorizer)

	connectionClient := botservice.NewBotConnectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&connectionClient.Client, o.ResourceManagerAuthorizer)

	channelClient := botservice.NewChannelsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&channelClient.Client, o.ResourceManagerAuthorizer)

	healthBotsClient := healthbot_2022_08_08.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	return &Client{
		BotClient:        &botClient,
		ChannelClient:    &channelClient,
		ConnectionClient: &connectionClient,
		HealthBotClient:  &healthBotsClient,
	}
}
