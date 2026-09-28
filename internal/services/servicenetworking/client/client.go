// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	v2025_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01" //azignore:AZG010 - package name does not match its path
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*v2025_01_01.Client, error) {
	client, err := v2025_01_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building ServiceNetworking client: %+v", err)
	}
	return client, nil
}
