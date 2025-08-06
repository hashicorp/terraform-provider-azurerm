// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/links"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	LinksClient         *links.LinksClient
	ServiceLinkerClient *servicelinker.ServicelinkerClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	linksClient, err := links.NewLinksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Links Client: %+v", err)
	}
	o.Configure(linksClient.Client, o.Authorizers.ResourceManager)

	serviceLinkerClient, err := servicelinker.NewServicelinkerClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ServiceLinker Client: %+v", err)
	}
	o.Configure(serviceLinkerClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		LinksClient:         linksClient,
		ServiceLinkerClient: serviceLinkerClient,
	}, nil
}
