// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2021-11-30/dedicatedhsms"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2025-03-31/cloudhsmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CloudHsmClustersClient *cloudhsmclusters.CloudHsmClustersClient
	DedicatedHsmClient     *dedicatedhsms.DedicatedHsmsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	cloudHsmClustersClient, err := cloudhsmclusters.NewCloudHsmClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building CloudHsmClusters client: %+v", err)
	}
	o.Configure(cloudHsmClustersClient.Client, o.Authorizers.ResourceManager)

	dedicatedHsmClient, err := dedicatedhsms.NewDedicatedHsmsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DedicatedHsms client: %+v", err)
	}
	o.Configure(dedicatedHsmClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		CloudHsmClustersClient: cloudHsmClustersClient,
		DedicatedHsmClient:     dedicatedHsmClient,
	}, nil
}
