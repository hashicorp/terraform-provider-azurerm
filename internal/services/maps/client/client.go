// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountsClient *accounts.AccountsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	accountsClient, err := accounts.NewAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, err
	}
	o.Configure(accountsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AccountsClient: accountsClient,
	}, nil
}
