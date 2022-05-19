package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/sdk/2021-05-13-preview/confidentialledger"
)

type Client struct {
	ConfidentialLedgereClient *confidentialledger.ConfidentialLedgerClient

	options *common.ClientOptions
}

func NewClient(o *common.ClientOptions) *Client {
	confidentialLedgerClient := confidentialledger.NewConfidentialLedgerClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&confidentialLedgerClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfidentialLedgereClient: &confidentialLedgerClient,
		options:                   o,
	}
}
