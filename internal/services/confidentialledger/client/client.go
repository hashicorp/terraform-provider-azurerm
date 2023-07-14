// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConfidentialLedgerClient *confidentialledger.ConfidentialLedgerClient

	options *common.ClientOptions
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	confidentialLedgerClient, err := confidentialledger.NewConfidentialLedgerClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConfidentialLedger client: %+v", err)
	}
	o.Configure(confidentialLedgerClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ConfidentialLedgerClient: confidentialLedgerClient,
		options:                  o,
	}, nil
}
