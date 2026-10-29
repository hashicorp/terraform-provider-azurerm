// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2025-07-25/openshiftclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/externalauths"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/hcpopenshiftclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/nodepools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	OpenShiftClustersClient    *openshiftclusters.OpenShiftClustersClient
	HcpOpenShiftClustersClient *hcpopenshiftclusters.HcpOpenShiftClustersClient
	ExternalAuthsClient        *externalauths.ExternalAuthsClient
	NodePoolsClient            *nodepools.NodePoolsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	openShiftClustersClient, err := openshiftclusters.NewOpenShiftClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("instantiating OpenShiftClustersClient: %+v", err)
	}
	o.Configure(openShiftClustersClient.Client, o.Authorizers.ResourceManager)

	hcpOpenShiftClustersClient, err := hcpopenshiftclusters.NewHcpOpenShiftClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("instantiating HcpOpenShiftClustersClient: %+v", err)
	}
	o.Configure(hcpOpenShiftClustersClient.Client, o.Authorizers.ResourceManager)

	externalAuthsClient, err := externalauths.NewExternalAuthsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExternalAuthsClient: %+v", err)
	}
	o.Configure(externalAuthsClient.Client, o.Authorizers.ResourceManager)

	nodePoolsClient, err := nodepools.NewNodePoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("instantiating NodePoolsClient: %+v", err)
	}
	o.Configure(nodePoolsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		OpenShiftClustersClient:    openShiftClustersClient,
		HcpOpenShiftClustersClient: hcpOpenShiftClustersClient,
		ExternalAuthsClient:        externalAuthsClient,
		NodePoolsClient:            nodePoolsClient,
	}, nil
}
