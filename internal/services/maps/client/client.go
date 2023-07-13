// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2021-02-01/accounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2021-02-01/creators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountsClient *accounts.AccountsClient
	CreatorsClient *creators.CreatorsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	accountsClient, err := accounts.NewAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, err
	}
	o.Configure(accountsClient.Client, o.Authorizers.ResourceManager)

	creatorsClient, err := creators.NewCreatorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, err
	}
	o.Configure(creatorsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AccountsClient: accountsClient,
		CreatorsClient: creatorsClient,
	}, nil
}
