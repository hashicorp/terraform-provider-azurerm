// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2021-10-01/connectedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ArcKubernetesClient     *connectedclusters.ConnectedClustersClient
	ExtensionsClient        *extensions.ExtensionsClient
	FluxConfigurationClient *fluxconfiguration.FluxConfigurationClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	arcKubernetesClient, err := connectedclusters.NewConnectedClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConnectedClusters client: %+v", err)
	}
	o.Configure(arcKubernetesClient.Client, o.Authorizers.ResourceManager)

	extensionsClient, err := extensions.NewExtensionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Extensions client: %+v", err)
	}
	o.Configure(extensionsClient.Client, o.Authorizers.ResourceManager)

	fluxConfigurationClient, err := fluxconfiguration.NewFluxConfigurationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FluxConfiguration client: %+v", err)
	}
	o.Configure(fluxConfigurationClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ArcKubernetesClient:     arcKubernetesClient,
		ExtensionsClient:        extensionsClient,
		FluxConfigurationClient: fluxConfigurationClient,
	}, nil
}
