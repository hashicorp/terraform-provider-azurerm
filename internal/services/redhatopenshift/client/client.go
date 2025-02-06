// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2023-09-04/openshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	OpenShiftClustersClient *openshiftclusters.OpenShiftClustersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	openShiftClustersClient, err := openshiftclusters.NewOpenShiftClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("instantiating OpenShiftClustersClient: %+v", err)
	}
	o.Configure(openShiftClustersClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		OpenShiftClustersClient: openShiftClustersClient,
	}, nil
}
