package client

import (
	"github.com/Azure/azure-sdk-for-go/services/policyinsights/mgmt/2019-10-01/policyinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	RemediationsClient *policyinsights.RemediationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	remediationsClient := policyinsights.NewRemediationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&remediationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RemediationsClient: &remediationsClient,
	}
}
