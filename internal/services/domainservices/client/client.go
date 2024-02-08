// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DomainServicesClient *domainservices.DomainServicesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	domainServicesClient, err := domainservices.NewDomainServicesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DomainServices client: %+v", err)
	}
	o.Configure(domainServicesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		DomainServicesClient: domainServicesClient,
	}, nil
}
