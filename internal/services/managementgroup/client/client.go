// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/managementgroups/2020-05-01/managementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	GroupsClient *managementgroups.ManagementGroupsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	groupsClient, err := managementgroups.NewManagementGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building managementgroups client: %+v", err)
	}
	o.Configure(groupsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		GroupsClient: groupsClient,
	}, nil
}
