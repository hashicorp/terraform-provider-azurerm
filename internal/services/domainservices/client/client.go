// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DomainServicesClient *domainservices.DomainServicesClient
}

func NewClient(o *common.ClientOptions) *Client {
	domainServicesClient := domainservices.NewDomainServicesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&domainServicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DomainServicesClient: &domainServicesClient,
	}
}
