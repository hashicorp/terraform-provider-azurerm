// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storageactions/2023-01-01/storagetasks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	StorageTasksClient *storagetasks.StorageTasksClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	storageTasksClient, err := storagetasks.NewStorageTasksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Storage Tasks Client: %+v", err)
	}
	o.Configure(storageTasksClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		StorageTasksClient: storageTasksClient,
	}, nil
}
