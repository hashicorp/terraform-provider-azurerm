// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2018-06-01/hdinsight" // nolint: staticcheck
	hdinsight_v2021_06_01 "github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*hdinsight_v2021_06_01.Client

	ClustersClient *hdinsight.ClustersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	// due to a bug in the HDInsight API we can't reuse client with the same x-ms-correlation-request-id for multiple updates
	opts := *o
	opts.DisableCorrelationRequestID = true

	client, err := hdinsight_v2021_06_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)

		// due to a bug in the HDInsight API we can't reuse client with the same x-ms-correlation-request-id for multiple updates
		c.CorrelationId = ""
	})
	if err != nil {
		return nil, fmt.Errorf("building meta client: %+v", err)
	}

	ClustersClient := hdinsight.NewClustersClientWithBaseURI(opts.ResourceManagerEndpoint, opts.SubscriptionId)
	opts.ConfigureClient(&ClustersClient.Client, opts.ResourceManagerAuthorizer)

	return &Client{
		Client: client,

		ClustersClient: &ClustersClient,
	}, nil
}
