package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/sdk/2020-05-01/frontdoors"
)

type Client struct {
	FrontDoorsClient *frontdoors.FrontDoorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	frontDoorsClient := frontdoors.NewFrontDoorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FrontDoorsClient: &frontDoorsClient,
	}
}
