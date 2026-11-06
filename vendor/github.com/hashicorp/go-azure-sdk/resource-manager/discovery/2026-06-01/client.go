package v2026_06_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelfprivateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelfprivatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/chatmodeldeployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/nodepools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/storageassets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/storagecontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/supercomputers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/tools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/workspaceprivateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/workspaceprivatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	BookshelfPrivateEndpointConnections *bookshelfprivateendpointconnections.BookshelfPrivateEndpointConnectionsClient
	BookshelfPrivateLinkResources       *bookshelfprivatelinkresources.BookshelfPrivateLinkResourcesClient
	Bookshelves                         *bookshelves.BookshelvesClient
	ChatModelDeployments                *chatmodeldeployments.ChatModelDeploymentsClient
	NodePools                           *nodepools.NodePoolsClient
	Projects                            *projects.ProjectsClient
	StorageAssets                       *storageassets.StorageAssetsClient
	StorageContainers                   *storagecontainers.StorageContainersClient
	Supercomputers                      *supercomputers.SupercomputersClient
	Tools                               *tools.ToolsClient
	WorkspacePrivateEndpointConnections *workspaceprivateendpointconnections.WorkspacePrivateEndpointConnectionsClient
	WorkspacePrivateLinkResources       *workspaceprivatelinkresources.WorkspacePrivateLinkResourcesClient
	Workspaces                          *workspaces.WorkspacesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	bookshelfPrivateEndpointConnectionsClient, err := bookshelfprivateendpointconnections.NewBookshelfPrivateEndpointConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building BookshelfPrivateEndpointConnections client: %+v", err)
	}
	configureFunc(bookshelfPrivateEndpointConnectionsClient.Client)

	bookshelfPrivateLinkResourcesClient, err := bookshelfprivatelinkresources.NewBookshelfPrivateLinkResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building BookshelfPrivateLinkResources client: %+v", err)
	}
	configureFunc(bookshelfPrivateLinkResourcesClient.Client)

	bookshelvesClient, err := bookshelves.NewBookshelvesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Bookshelves client: %+v", err)
	}
	configureFunc(bookshelvesClient.Client)

	chatModelDeploymentsClient, err := chatmodeldeployments.NewChatModelDeploymentsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ChatModelDeployments client: %+v", err)
	}
	configureFunc(chatModelDeploymentsClient.Client)

	nodePoolsClient, err := nodepools.NewNodePoolsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NodePools client: %+v", err)
	}
	configureFunc(nodePoolsClient.Client)

	projectsClient, err := projects.NewProjectsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Projects client: %+v", err)
	}
	configureFunc(projectsClient.Client)

	storageAssetsClient, err := storageassets.NewStorageAssetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StorageAssets client: %+v", err)
	}
	configureFunc(storageAssetsClient.Client)

	storageContainersClient, err := storagecontainers.NewStorageContainersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StorageContainers client: %+v", err)
	}
	configureFunc(storageContainersClient.Client)

	supercomputersClient, err := supercomputers.NewSupercomputersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Supercomputers client: %+v", err)
	}
	configureFunc(supercomputersClient.Client)

	toolsClient, err := tools.NewToolsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Tools client: %+v", err)
	}
	configureFunc(toolsClient.Client)

	workspacePrivateEndpointConnectionsClient, err := workspaceprivateendpointconnections.NewWorkspacePrivateEndpointConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building WorkspacePrivateEndpointConnections client: %+v", err)
	}
	configureFunc(workspacePrivateEndpointConnectionsClient.Client)

	workspacePrivateLinkResourcesClient, err := workspaceprivatelinkresources.NewWorkspacePrivateLinkResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building WorkspacePrivateLinkResources client: %+v", err)
	}
	configureFunc(workspacePrivateLinkResourcesClient.Client)

	workspacesClient, err := workspaces.NewWorkspacesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Workspaces client: %+v", err)
	}
	configureFunc(workspacesClient.Client)

	return &Client{
		BookshelfPrivateEndpointConnections: bookshelfPrivateEndpointConnectionsClient,
		BookshelfPrivateLinkResources:       bookshelfPrivateLinkResourcesClient,
		Bookshelves:                         bookshelvesClient,
		ChatModelDeployments:                chatModelDeploymentsClient,
		NodePools:                           nodePoolsClient,
		Projects:                            projectsClient,
		StorageAssets:                       storageAssetsClient,
		StorageContainers:                   storageContainersClient,
		Supercomputers:                      supercomputersClient,
		Tools:                               toolsClient,
		WorkspacePrivateEndpointConnections: workspacePrivateEndpointConnectionsClient,
		WorkspacePrivateLinkResources:       workspacePrivateLinkResourcesClient,
		Workspaces:                          workspacesClient,
	}, nil
}
