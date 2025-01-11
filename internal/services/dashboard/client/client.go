// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2023-09-01/grafanaresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2023-09-01/managedprivateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	GrafanaResourceClient         *grafanaresource.GrafanaResourceClient
	ManagedPrivateEndpointsClient *managedprivateendpoints.ManagedPrivateEndpointsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	grafanaResourceClient, err := grafanaresource.NewGrafanaResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GrafanaResource client: %+v", err)
	}
	o.Configure(grafanaResourceClient.Client, o.Authorizers.ResourceManager)

	managedPrivateEndpointsClient, err := managedprivateendpoints.NewManagedPrivateEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ManagedPrivateEndpoints client: %+v", err)
	}

	o.Configure(managedPrivateEndpointsClient.Client, o.Authorizers.ResourceManager)
	return &Client{
		GrafanaResourceClient:         grafanaResourceClient,
		ManagedPrivateEndpointsClient: managedPrivateEndpointsClient,
	}, nil
}
