// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypackqueries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypacks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/dataexport"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedstorageaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/storageinsights"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	featureWorkspaces "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/deletedworkspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationsmanagement/2015-11-01-preview/solution"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ClusterClient              *clusters.ClustersClient
	DataExportClient           *dataexport.DataExportClient
	DataSourcesClient          *datasources.DataSourcesClient
	DeletedWorkspacesClient    *deletedworkspaces.DeletedWorkspacesClient
	LinkedServicesClient       *linkedservices.LinkedServicesClient
	LinkedStorageAccountClient *linkedstorageaccounts.LinkedStorageAccountsClient
	QueryPacksClient           *querypacks.QueryPacksClient
	SavedSearchesClient        *savedsearches.SavedSearchesClient
	SolutionsClient            *solution.SolutionClient
	StorageInsightsClient      *storageinsights.StorageInsightsClient
	QueryPackQueriesClient     *querypackqueries.QueryPackQueriesClient
	SharedKeyWorkspacesClient  *workspaces.WorkspacesClient
	TablesClient               *tables.TablesClient
	WorkspaceClient            *featureWorkspaces.WorkspacesClient // 2022-10-01 API version does not contain sharedkeys related API, so we keep two versions SDK of this API
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	clusterClient, err := clusters.NewClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Clusters client: %+v", err)
	}
	o.Configure(clusterClient.Client, o.Authorizers.ResourceManager)

	dataExportClient, err := dataexport.NewDataExportClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DataExport client: %+v", err)
	}
	o.Configure(dataExportClient.Client, o.Authorizers.ResourceManager)

	dataSourcesClient, err := datasources.NewDataSourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DataSources client: %+v", err)
	}
	o.Configure(dataSourcesClient.Client, o.Authorizers.ResourceManager)

	deletedWorkspacesClient, err := deletedworkspaces.NewDeletedWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Deleted Workspaces Client: %+v", err)
	}
	o.Configure(deletedWorkspacesClient.Client, o.Authorizers.ResourceManager)

	workspacesClient, err := workspaces.NewWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Workspaces client: %+v", err)
	}
	o.Configure(workspacesClient.Client, o.Authorizers.ResourceManager)

	featureWorkspaceClient, err := featureWorkspaces.NewWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FeatureWorkspaces client: %+v", err)
	}
	o.Configure(featureWorkspaceClient.Client, o.Authorizers.ResourceManager)

	savedSearchesClient, err := savedsearches.NewSavedSearchesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SavedSearches client: %+v", err)
	}
	o.Configure(savedSearchesClient.Client, o.Authorizers.ResourceManager)

	solutionsClient, err := solution.NewSolutionClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Solutions client: %+v", err)
	}
	o.Configure(solutionsClient.Client, o.Authorizers.ResourceManager)

	storageInsightsClient, err := storageinsights.NewStorageInsightsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building StorageInsights client: %+v", err)
	}
	o.Configure(storageInsightsClient.Client, o.Authorizers.ResourceManager)

	linkedServicesClient, err := linkedservices.NewLinkedServicesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building LinkedServices client: %+v", err)
	}
	o.Configure(linkedServicesClient.Client, o.Authorizers.ResourceManager)

	linkedStorageAccountClient, err := linkedstorageaccounts.NewLinkedStorageAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building LinkedStorageAccounts client: %+v", err)
	}
	o.Configure(linkedStorageAccountClient.Client, o.Authorizers.ResourceManager)

	queryPacksClient, err := querypacks.NewQueryPacksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building QueryPacks client: %+v", err)
	}
	o.Configure(queryPacksClient.Client, o.Authorizers.ResourceManager)

	queryPackQueriesClient, err := querypackqueries.NewQueryPackQueriesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building QueryPackQueries client: %+v", err)
	}
	o.Configure(queryPackQueriesClient.Client, o.Authorizers.ResourceManager)

	tablesClient, err := tables.NewTablesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Tables client: %+v", err)
	}
	o.Configure(tablesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ClusterClient:              clusterClient,
		DataExportClient:           dataExportClient,
		DataSourcesClient:          dataSourcesClient,
		DeletedWorkspacesClient:    deletedWorkspacesClient,
		LinkedServicesClient:       linkedServicesClient,
		LinkedStorageAccountClient: linkedStorageAccountClient,
		QueryPacksClient:           queryPacksClient,
		QueryPackQueriesClient:     queryPackQueriesClient,
		SavedSearchesClient:        savedSearchesClient,
		SolutionsClient:            solutionsClient,
		StorageInsightsClient:      storageInsightsClient,
		SharedKeyWorkspacesClient:  workspacesClient,
		TablesClient:               tablesClient,
		WorkspaceClient:            featureWorkspaceClient,
	}, nil
}
