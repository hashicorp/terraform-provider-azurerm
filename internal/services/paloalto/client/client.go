// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	paloalto_2025_10_08 "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*paloalto_2025_10_08.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := paloalto_2025_10_08.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for PaloAlto: %+v", err)
	}

	return &Client{
		Client: client,
	}, nil
}
