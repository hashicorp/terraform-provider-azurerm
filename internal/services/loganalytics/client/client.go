// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypackqueries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypacks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/dataexport"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedstorageaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/storageinsights"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	featureWorkspaces "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationsmanagement/2015-11-01-preview/solution"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ClusterClient              *clusters.ClustersClient
	DataExportClient           *dataexport.DataExportClient
	DataSourcesClient          *datasources.DataSourcesClient
	LinkedServicesClient       *linkedservices.LinkedServicesClient
	LinkedStorageAccountClient *linkedstorageaccounts.LinkedStorageAccountsClient
	QueryPacksClient           *querypacks.QueryPacksClient
	SavedSearchesClient        *savedsearches.SavedSearchesClient
	SolutionsClient            *solution.SolutionClient
	StorageInsightsClient      *storageinsights.StorageInsightsClient
	QueryPackQueriesClient     *querypackqueries.QueryPackQueriesClient
	SharedKeyWorkspacesClient  *workspaces.WorkspacesClient
	WorkspaceClient            *featureWorkspaces.WorkspacesClient // 2022-10-01 API version does not contain sharedkeys related API, so we keep two versions SDK of this API
}

func NewClient(o *common.ClientOptions) *Client {
	ClusterClient := clusters.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ClusterClient.Client, o.ResourceManagerAuthorizer)

	DataExportClient := dataexport.NewDataExportClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DataExportClient.Client, o.ResourceManagerAuthorizer)

	DataSourcesClient := datasources.NewDataSourcesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DataSourcesClient.Client, o.ResourceManagerAuthorizer)

	WorkspacesClient := workspaces.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	featureWorkspaceClient := featureWorkspaces.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&featureWorkspaceClient.Client, o.ResourceManagerAuthorizer)

	SavedSearchesClient := savedsearches.NewSavedSearchesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&SavedSearchesClient.Client, o.ResourceManagerAuthorizer)

	SolutionsClient := solution.NewSolutionClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&SolutionsClient.Client, o.ResourceManagerAuthorizer)

	StorageInsightsClient := storageinsights.NewStorageInsightsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&StorageInsightsClient.Client, o.ResourceManagerAuthorizer)

	LinkedServicesClient := linkedservices.NewLinkedServicesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&LinkedServicesClient.Client, o.ResourceManagerAuthorizer)

	LinkedStorageAccountClient := linkedstorageaccounts.NewLinkedStorageAccountsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&LinkedStorageAccountClient.Client, o.ResourceManagerAuthorizer)

	QueryPacksClient := querypacks.NewQueryPacksClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&QueryPacksClient.Client, o.ResourceManagerAuthorizer)

	QueryPackQueriesClient := querypackqueries.NewQueryPackQueriesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&QueryPackQueriesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ClusterClient:              &ClusterClient,
		DataExportClient:           &DataExportClient,
		DataSourcesClient:          &DataSourcesClient,
		LinkedServicesClient:       &LinkedServicesClient,
		LinkedStorageAccountClient: &LinkedStorageAccountClient,
		QueryPacksClient:           &QueryPacksClient,
		QueryPackQueriesClient:     &QueryPackQueriesClient,
		SavedSearchesClient:        &SavedSearchesClient,
		SolutionsClient:            &SolutionsClient,
		StorageInsightsClient:      &StorageInsightsClient,
		SharedKeyWorkspacesClient:  &WorkspacesClient,
		WorkspaceClient:            &featureWorkspaceClient,
	}
}
