// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/videoanalyzer/2021-05-01-preview/edgemodules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/videoanalyzer/2021-05-01-preview/videoanalyzers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	EdgeModuleClient     *edgemodules.EdgeModulesClient
	VideoAnalyzersClient *videoanalyzers.VideoAnalyzersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	edgeModulesClient, err := edgemodules.NewEdgeModulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Edge Modules Client: %+v", err)
	}
	o.Configure(edgeModulesClient.Client, o.Authorizers.ResourceManager)

	videoAnalyzersClient, err := videoanalyzers.NewVideoAnalyzersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Video Analyzers Client: %+v", err)
	}
	o.Configure(videoAnalyzersClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		EdgeModuleClient:     edgeModulesClient,
		VideoAnalyzersClient: videoAnalyzersClient,
	}, nil
}
