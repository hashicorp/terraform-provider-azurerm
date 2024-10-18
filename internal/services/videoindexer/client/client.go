package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/videoindexer/2024-01-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountClient *accounts.AccountsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	accountClient, err := accounts.NewAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Video Indexer Account client: %+v", err)
	}
	o.Configure(accountClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AccountClient: accountClient,
	}, nil
}
