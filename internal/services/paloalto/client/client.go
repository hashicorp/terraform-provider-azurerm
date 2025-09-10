// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	paloalto_2022_08_29 "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29"
	paloalto_2025_05_23 "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-05-23"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*paloalto_2022_08_29.Client
	PaloAltoClient_v2025_05_23 *paloalto_2025_05_23.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := paloalto_2022_08_29.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Network: %+v", err)
	}

	paloAltoClient_v2025_05_23, err := paloalto_2025_05_23.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:                     client,
		PaloAltoClient_v2025_05_23: paloAltoClient_v2025_05_23,
	}, nil
}
