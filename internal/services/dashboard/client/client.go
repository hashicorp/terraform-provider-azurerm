// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedgrafanas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedprivateendpointmodels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	GrafanaResourceClient         *managedgrafanas.ManagedGrafanasClient
	ManagedPrivateEndpointsClient *managedprivateendpointmodels.ManagedPrivateEndpointModelsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	grafanaResourceClient, err := managedgrafanas.NewManagedGrafanasClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GrafanaResource client: %+v", err)
	}
	o.Configure(grafanaResourceClient.Client, o.Authorizers.ResourceManager)

	managedPrivateEndpointsClient, err := managedprivateendpointmodels.NewManagedPrivateEndpointModelsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ManagedPrivateEndpoints client: %+v", err)
	}

	o.Configure(managedPrivateEndpointsClient.Client, o.Authorizers.ResourceManager)
	return &Client{
		GrafanaResourceClient:         grafanaResourceClient,
		ManagedPrivateEndpointsClient: managedPrivateEndpointsClient,
	}, nil
}
