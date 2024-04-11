// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	emailchannel_2022_09_15 "github.com/hashicorp/go-azure-sdk/resource-manager/botservice/2022-09-15/channel"
	healthbot_2022_08_08 "github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2022-08-08"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

type Client struct {
	BotClient          *botservice.BotsClient
	ConnectionClient   *botservice.BotConnectionClient
	ChannelClient      *botservice.ChannelsClient
	EmailChannelClient *emailchannel_2022_09_15.ChannelClient
	HealthBotClient    *healthbot_2022_08_08.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	botClient := botservice.NewBotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&botClient.Client, o.ResourceManagerAuthorizer)

	connectionClient := botservice.NewBotConnectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&connectionClient.Client, o.ResourceManagerAuthorizer)

	channelClient := botservice.NewChannelsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&channelClient.Client, o.ResourceManagerAuthorizer)

	emailChannelClient, err := emailchannel_2022_09_15.NewChannelClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building EmailChannels client: %+v", err)
	}
	o.Configure(emailChannelClient.Client, o.Authorizers.ResourceManager)

	healthBotsClient, err := healthbot_2022_08_08.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
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
