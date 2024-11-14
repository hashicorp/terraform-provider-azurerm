// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/deployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/metadataschemas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ApiDefinitionsClient  *apidefinitions.ApiDefinitionsClient
	ApisClient            *apis.ApisClient
	ApiVersionsClient     *apiversions.ApiVersionsClient
	DeploymentsClient     *deployments.DeploymentsClient
	EnvironmentsClient    *environments.EnvironmentsClient
	MetadataSchemasClient *metadataschemas.MetadataSchemasClient
	ServicesClient        *services.ServicesClient
	WorkspacesClient      *workspaces.WorkspacesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	apiDefinitionsClient, err := apidefinitions.NewApiDefinitionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Definitions client: %+v", err)
	}
	o.Configure(apiDefinitionsClient.Client, o.Authorizers.ResourceManager)

	apisClient, err := apis.NewApisClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Apis client: %+v", err)
	}
	o.Configure(apisClient.Client, o.Authorizers.ResourceManager)

	apiVersionsClient, err := apiversions.NewApiVersionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Api Versions client: %+v", err)
	}
	o.Configure(apiVersionsClient.Client, o.Authorizers.ResourceManager)

	deploymentsClient, err := deployments.NewDeploymentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Deployments client: %+v", err)
	}
	o.Configure(deploymentsClient.Client, o.Authorizers.ResourceManager)

	environmentsClient, err := environments.NewEnvironmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Environments client: %+v", err)
	}
	o.Configure(environmentsClient.Client, o.Authorizers.ResourceManager)

	metadataSchemasClient, err := metadataschemas.NewMetadataSchemasClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building MetadataSchemas client: %+v", err)
	}
	o.Configure(metadataSchemasClient.Client, o.Authorizers.ResourceManager)

	servicesClient, err := services.NewServicesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Services client: %+v", err)
	}
	o.Configure(servicesClient.Client, o.Authorizers.ResourceManager)

	workspacesClient, err := workspaces.NewWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Workspaces client: %+v", err)
	}
	o.Configure(workspacesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ApiDefinitionsClient:  apiDefinitionsClient,
		ApisClient:            apisClient,
		ApiVersionsClient:     apiVersionsClient,
		DeploymentsClient:     deploymentsClient,
		EnvironmentsClient:    environmentsClient,
		MetadataSchemasClient: metadataSchemasClient,
		ServicesClient:        servicesClient,
		WorkspacesClient:      workspacesClient,
	}, nil
}
