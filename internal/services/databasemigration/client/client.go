// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/projectresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/serviceresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ServicesClient *serviceresource.ServiceResourceClient
	ProjectsClient *projectresource.ProjectResourceClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	servicesClient, err := serviceresource.NewServiceResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ServicesClient client: %+v", err)
	}
	o.Configure(servicesClient.Client, o.Authorizers.ResourceManager)

	projectsClient, err := projectresource.NewProjectResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ProjectsClient client: %+v", err)
	}
	o.Configure(projectsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ServicesClient: servicesClient,
		ProjectsClient: projectsClient,
	}, nil
}
