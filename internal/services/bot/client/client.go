// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/botservice/2022-09-15/channel"
	v2025_05_25 "github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2025-05-25" //azignore:AZG010 - package name does not match its path
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

type Client struct {
	BotClient          *botservice.BotsClient
	ConnectionClient   *botservice.BotConnectionClient
	ChannelClient      *botservice.ChannelsClient
	EmailChannelClient *channel.ChannelClient
	HealthBotClient    *v2025_05_25.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	botClient := botservice.NewBotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&botClient.Client, o.ResourceManagerAuthorizer)

	connectionClient := botservice.NewBotConnectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&connectionClient.Client, o.ResourceManagerAuthorizer)

	channelClient := botservice.NewChannelsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&channelClient.Client, o.ResourceManagerAuthorizer)

	emailChannelClient, err := channel.NewChannelClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building EmailChannels client: %+v", err)
	}
	o.Configure(emailChannelClient.Client, o.Authorizers.ResourceManager)

	healthBotsClient, err := v2025_05_25.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building HealthBots client: %+v", err)
	}

	return &Client{
		BotClient:          &botClient,
		ChannelClient:      &channelClient,
		ConnectionClient:   &connectionClient,
		EmailChannelClient: emailChannelClient,
		HealthBotClient:    healthBotsClient,
	}, nil
}
