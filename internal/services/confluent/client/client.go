// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/organizationresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	OrganizationResourcesClient *organizationresources.OrganizationResourcesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	organizationResourcesClient, err := organizationresources.NewOrganizationResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building OrganizationResources client: %+v", err)
	}
	o.Configure(organizationResourcesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		OrganizationResourcesClient: organizationResourcesClient,
	}, nil
}
