// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	fluidrelay_2022_05_26 "github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*fluidrelay_2022_05_26.Client, error) {
	client, err := fluidrelay_2022_05_26.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building FluidRelay client: %+v", err)
	}

	return client, nil
}
