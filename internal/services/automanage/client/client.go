// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

type Client struct {
	ConfigurationClient *automanage.ConfigurationProfilesClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationProfileClient := automanage.NewConfigurationProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfileClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationClient: &configurationProfileClient,
	}
}
