// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	v2018_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01" //azignore:AZG010 - package name does not match its path
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*v2018_05_01.Client, error) {
	client, err := v2018_05_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
