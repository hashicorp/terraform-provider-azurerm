package client

import (
	"github.com/Azure/azure-sdk-for-go/services/healthbot/mgmt/2020-12-08/healthbot" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest"
	healthbot_2020_12_08 "github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2020-12-08"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

type Client struct {
	BotClient        *botservice.BotsClient
	ConnectionClient *botservice.BotConnectionClient
	ChannelClient    *botservice.ChannelsClient
	HealthBotClient  *healthbot_2020_12_08.Client
}

func NewClient(o *common.ClientOptions) *Client {
	botClient := botservice.NewBotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&botClient.Client, o.ResourceManagerAuthorizer)

	connectionClient := botservice.NewBotConnectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&connectionClient.Client, o.ResourceManagerAuthorizer)

	channelClient := botservice.NewChannelsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&channelClient.Client, o.ResourceManagerAuthorizer)

	healthBotClient := healthbot.NewBotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&healthBotClient.Client, o.ResourceManagerAuthorizer)

	healthBotsClient := healthbot_2020_12_08.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	return &Client{
		BotClient:        &botClient,
		ChannelClient:    &channelClient,
		ConnectionClient: &connectionClient,
		HealthBotClient:  &healthBotsClient,
	}
}
