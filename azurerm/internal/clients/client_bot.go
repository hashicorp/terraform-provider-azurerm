package clients

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type BotClient struct {
	BotClient        *botservice.BotsClient
	ConnectionClient *botservice.BotConnectionClient
	ChannelClient    *botservice.ChannelsClient
}

func newBotClient(o *common.ClientOptions) *BotClient {
	botClient := botservice.NewBotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&botClient.Client, o.ResourceManagerAuthorizer)

	connectionClient := botservice.NewBotConnectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&connectionClient.Client, o.ResourceManagerAuthorizer)

	channelClient := botservice.NewChannelsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&channelClient.Client, o.ResourceManagerAuthorizer)

	return &BotClient{
		BotClient:        &botClient,
		ChannelClient:    &channelClient,
		ConnectionClient: &connectionClient,
	}
}
