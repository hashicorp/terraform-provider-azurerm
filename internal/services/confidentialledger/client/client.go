package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConfidentialLedgerClient *confidentialledger.ConfidentialLedgerClient

	options *common.ClientOptions
}

func NewClient(o *common.ClientOptions) *Client {
	confidentialLedgerClient := confidentialledger.NewConfidentialLedgerClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&confidentialLedgerClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfidentialLedgerClient: &confidentialLedgerClient,
		options:                  o,
	}
}
