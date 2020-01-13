package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/managementpartner/mgmt/2018-02-01/managementpartner"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	PartnerClient *managementpartner.PartnerClient
}

func NewClient(o *common.ClientOptions) *Client {
	partnerClient := managementpartner.NewPartnerClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&partnerClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		PartnerClient: &partnerClient,
	}
}
