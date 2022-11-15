package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/lab"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	LabClient *lab.LabClient
}

func NewClient(o *common.ClientOptions) *Client {
	LabClient := lab.NewLabClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&LabClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		LabClient: &LabClient,
	}
}
