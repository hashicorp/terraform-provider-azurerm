// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	loadtestserviceV20221201 "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2025-09-01/playwrightworkspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20221201                  loadtestserviceV20221201.Client
	PlaywrightWorkspacesClient *playwrightworkspaces.PlaywrightWorkspacesClient
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {
	v20221201Client, err := loadtestserviceV20221201.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for loadtestservice V20221201: %+v", err)
	}

	playwrightWorkspacesClient, err := playwrightworkspaces.NewPlaywrightWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PlaywrightWorkspaces Client: %+v", err)
	}
	o.Configure(playwrightWorkspacesClient.Client, o.Authorizers.ResourceManager)

	return &AutoClient{
		V20221201:                  *v20221201Client,
		PlaywrightWorkspacesClient: playwrightWorkspacesClient,
	}, nil
}
