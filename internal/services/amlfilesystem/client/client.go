package client

import (
	"fmt"

	amlfilesystem_2023_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*amlfilesystem_2023_05_01.Client, error) {
	client, err := amlfilesystem_2023_05_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building AML File System client: %+v", err)
	}

	return client, nil
}
