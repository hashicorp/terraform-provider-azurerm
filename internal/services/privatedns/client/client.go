// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2020-06-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2020-06-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2020-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	RecordSetsClient          *recordsets.RecordSetsClient
	PrivateZonesClient        *privatezones.PrivateZonesClient
	VirtualNetworkLinksClient *virtualnetworklinks.VirtualNetworkLinksClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	recordSetsClient, err := recordsets.NewRecordSetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Record Sets Client: %v", err)
	}
	o.Configure(recordSetsClient.Client, o.Authorizers.ResourceManager)

	privateZonesClient, err := privatezones.NewPrivateZonesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Private Zones client: %v", err)
	}
	o.Configure(privateZonesClient.Client, o.Authorizers.ResourceManager)

	virtualNetworkLinksClient, err := virtualnetworklinks.NewVirtualNetworkLinksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Virtual Network Links client: %v", err)
	}
	o.Configure(virtualNetworkLinksClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		RecordSetsClient:          recordSetsClient,
		PrivateZonesClient:        privateZonesClient,
		VirtualNetworkLinksClient: virtualNetworkLinksClient,
	}, nil
}
