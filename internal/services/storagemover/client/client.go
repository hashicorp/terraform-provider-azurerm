package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	StorageMoversClient *storagemovers.StorageMoversClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	storageMoversClient, err := storagemovers.NewStorageMoversClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Storage Movers client: %+v", err)
	}
	o.Configure(storageMoversClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		StorageMoversClient: storageMoversClient,
	}, nil
}
