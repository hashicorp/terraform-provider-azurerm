// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	preflightvalidation "github.com/hashicorp/terraform-provider-azurerm/internal/preflight/sdk" //azignore:AZG010 - package name does not match its path
)

type Client struct {
	PreflightClient *preflightvalidation.PreflightClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	preflightClient, err := preflightvalidation.NewResourceValidationClientClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Preflight client: %+v", err)
	}
	o.Configure(preflightClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		PreflightClient: preflightClient,
	}, nil
}
