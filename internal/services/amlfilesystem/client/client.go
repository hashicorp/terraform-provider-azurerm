package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/amlfilesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AMLFileSystemClient *amlfilesystems.AmlFilesystemsClient
}

func NewClient(o *common.ClientOptions) *Client {
	amlFileSystemClient := amlfilesystems.NewAmlFilesystemsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&amlFileSystemClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AMLFileSystemClient: &amlFileSystemClient,
	}
}
