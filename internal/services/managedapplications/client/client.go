// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applicationdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ApplicationClient           *applications.ApplicationsClient
	ApplicationDefinitionClient *applicationdefinitions.ApplicationDefinitionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	applicationClient, err := applications.NewApplicationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Applications client: %+v", err)
	}
	o.Configure(applicationClient.Client, o.Authorizers.ResourceManager)

	applicationDefinitionClient, err := applicationdefinitions.NewApplicationDefinitionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Application Definitions client: %+v", err)
	}
	o.Configure(applicationDefinitionClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ApplicationClient:           applicationClient,
		ApplicationDefinitionClient: applicationDefinitionClient,
	}, nil
}
