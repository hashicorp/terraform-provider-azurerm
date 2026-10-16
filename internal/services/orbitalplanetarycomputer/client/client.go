// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/orbitalplanetarycomputer/2026-04-15/geocatalogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	GeoCatalogsClient *geocatalogs.GeoCatalogsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	geoCatalogsClient, err := geocatalogs.NewGeoCatalogsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GeoCatalogs Client: %+v", err)
	}
	o.Configure(geoCatalogsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		GeoCatalogsClient: geoCatalogsClient,
	}, nil
}
