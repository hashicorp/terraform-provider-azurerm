// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/qumulostorage/2024-06-19/filesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FileSystemsClient filesystems.FileSystemsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	fileSystemsClient, err := filesystems.NewFileSystemsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FileSystems client: %+v", err)
	}
	o.Configure(fileSystemsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		FileSystemsClient: *fileSystemsClient,
	}, nil
}
