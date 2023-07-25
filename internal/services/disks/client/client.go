// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/diskpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/iscsitargets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DiskPoolsClient            *diskpools.DiskPoolsClient
	DisksPoolIscsiTargetClient *iscsitargets.IscsiTargetsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	diskPoolsClient, err := diskpools.NewDiskPoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Disk Pools Client: %+v", err)
	}
	o.Configure(diskPoolsClient.Client, o.Authorizers.ResourceManager)

	iscsiTargetClient, err := iscsitargets.NewIscsiTargetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Iscsi Target Client: %+v", err)
	}
	o.Configure(iscsiTargetClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		DiskPoolsClient:            diskPoolsClient,
		DisksPoolIscsiTargetClient: iscsiTargetClient,
	}, nil
}
