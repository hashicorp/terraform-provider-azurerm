package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/sdk/2021-08-01/diskpools"
)

type Client struct {
	DiskPoolsClient *diskpools.DiskPoolsClient
}

func NewClient(o *common.ClientOptions) *Client {
	diskPoolsClient := diskpools.NewDiskPoolsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&diskPoolsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DiskPoolsClient: &diskPoolsClient,
	}
}
