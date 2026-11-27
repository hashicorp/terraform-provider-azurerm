// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2024-04-01/invoicesection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	InvoiceSectionClient *invoicesection.InvoiceSectionClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	invoiceSectionClient, err := invoicesection.NewInvoiceSectionClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Invoice Section client: %+v", err)
	}
	o.Configure(invoiceSectionClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		InvoiceSectionClient: invoiceSectionClient,
	}, nil
}
